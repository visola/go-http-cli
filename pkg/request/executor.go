package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/visola/go-http-cli/pkg/profile"
	"github.com/visola/go-http-cli/pkg/session"
	"github.com/visola/go-http-cli/pkg/util"
	"github.com/visola/variables/variables"
)

const defaultMaxRedirectCount = 10

// ExecuteRequestLoop executes HTTP requests based on the passed in options until there're no more
// requests to be executed.
func ExecuteRequestLoop(executionContext ExecutionContext) ([]ExecutedRequestResponse, error) {
	maxRedirectCount := util.FirstOrZero(executionContext.MaxRedirect, defaultMaxRedirectCount)
	client := createHTTPClient()

	mergedProfiles, profileError := profile.LoadAndMergeProfiles(executionContext.ProfileNames)
	if profileError != nil {
		return nil, profileError
	}

	initialVariables := mergeVariables(executionContext.Variables, mergedProfiles.Variables)

	requestsToExecute := []Request{executionContext.Request}
	result := make([]ExecutedRequestResponse, 0)
	redirectCount := 0
	for {
		currentConfiguredRequest := requestsToExecute[0]
		requestsToExecute = requestsToExecute[1:]

		var sessionErr error
		executionContext.Session, sessionErr = loadSessionForRequest(variables.ReplaceVariables(currentConfiguredRequest.URL, initialVariables))
		if sessionErr != nil {
			return nil, sessionErr
		}

		currentConfiguredRequest, processError := replaceRequestVariables(currentConfiguredRequest, mergedProfiles, executionContext)
		if processError != nil {
			return nil, processError
		}

		response, executeErr := executeRequest(client, currentConfiguredRequest, executionContext.Session)

		if executeErr != nil {
			return result, executeErr
		}

		requestResponse := ExecutedRequestResponse{
			Request:  currentConfiguredRequest,
			Response: *response,
		}
		result = append(result, requestResponse)

		postProcessOutput, postProcessError := PostProcess(&executionContext, result, executeErr)
		result[len(result)-1].PostProcessOutput = postProcessOutput
		if postProcessError != nil {
			result[len(result)-1].PostProcessError = postProcessError.Error()
		}

		if executeErr != nil {
			return result, executeErr
		}

		if shouldRedirect(response.StatusCode) && executionContext.FollowLocation == true {
			redirectCount++

			if redirectCount > maxRedirectCount {
				return result, fmt.Errorf("Max number of redirects reached: %d", maxRedirectCount)
			}

			redirectRequest := buildRedirect(&currentConfiguredRequest, response)
			requestsToExecute = append(requestsToExecute, *redirectRequest)
		}

		// If nothing else to execute, break
		if len(requestsToExecute) == 0 {
			break
		}
	}

	return result, nil
}

func buildRedirect(req *Request, response *Response) *Request {
	locationValues := response.Headers["Location"]
	if len(locationValues) > 0 && locationValues[0] != "" {
		location := locationValues[0]
		if !strings.HasPrefix(location, "http") {
			parsedURL, _ := url.Parse(req.URL)
			location = parsedURL.Scheme + "://" + parsedURL.Host + location
		}
		return &Request{
			URL: location,
		}
	}

	return nil
}

func createHTTPClient() *http.Client {
	return &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func executeRequest(client *http.Client, configuredRequest Request, currentSession *session.Session) (*Response, error) {
	httpRequest, httpRequestErr := BuildRequest(configuredRequest)
	if httpRequestErr != nil {
		return nil, httpRequestErr
	}

	httpResponse, httpResponseErr := client.Do(httpRequest)
	if httpResponseErr != nil {
		return nil, httpResponseErr
	}

	for _, cookie := range httpResponse.Cookies() {
		currentSession.Cookies[cookie.Name] = cookie
	}

	bodyBytes, readErr := ioutil.ReadAll(httpResponse.Body)

	if readErr != nil {
		return nil, readErr
	}

	headers := make(map[string][]string)
	for k, vs := range httpResponse.Header {
		headers[k] = append(headers[k], vs...)
	}

	return &Response{
		StatusCode: httpResponse.StatusCode,
		Status:     httpResponse.Status,
		Headers:    headers,
		Body:       string(bodyBytes),
		Protocol:   fmt.Sprintf("%d.%d", httpResponse.ProtoMajor, httpResponse.ProtoMinor),
	}, nil
}

func loadSessionForRequest(requestURL string) (*session.Session, error) {
	parsedURL, parseURLErr := url.Parse(requestURL)
	if parseURLErr != nil {
		return nil, parseURLErr
	}

	return session.Get(parsedURL.Hostname()), nil
}

func shouldRedirect(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently ||
		statusCode == http.StatusFound ||
		statusCode == http.StatusSeeOther
}
