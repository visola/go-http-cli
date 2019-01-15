package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
	"github.com/visola/go-http-cli/util"
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

	requestsToExecute := []Request{executionContext.Request}
	result := make([]ExecutedRequestResponse, 0)
	redirectCount := 0
	for {
		currentConfiguredRequest := requestsToExecute[0]
		requestsToExecute = requestsToExecute[1:]

		currentConfiguredRequest, processError := replaceRequestVariables(currentConfiguredRequest, mergedProfiles, executionContext.Variables)
		if processError != nil {
			return nil, processError
		}

		response, executeErr := executeRequest(client, currentConfiguredRequest)

		if executeErr != nil {
			return result, executeErr
		}

		requestResponse := ExecutedRequestResponse{
			Request:  currentConfiguredRequest,
			Response: *response,
		}
		result = append(result, requestResponse)

		postProcessOutput, postProcessError := PostProcess(executionContext.PostProcessCode, result, executeErr)
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

func executeRequest(client *http.Client, configuredRequest Request) (*Response, error) {
	httpRequest, httpRequestErr := BuildRequest(configuredRequest)
	if httpRequestErr != nil {
		return nil, httpRequestErr
	}

	httpResponse, httpResponseErr := client.Do(httpRequest)
	if httpResponseErr != nil {
		return nil, httpResponseErr
	}

	cookieErr := storeCookies(*httpRequest, *httpResponse)

	if cookieErr != nil {
		return nil, cookieErr
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

func shouldRedirect(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently ||
		statusCode == http.StatusFound ||
		statusCode == http.StatusSeeOther
}

func storeCookies(httpRequest http.Request, httpResponse http.Response) error {
	session, sessionErr := session.Get(httpRequest.URL.Hostname())
	if sessionErr != nil {
		return sessionErr
	}

	for _, cookie := range httpResponse.Cookies() {
		session.Cookies = append(session.Cookies, cookie)
	}

	return nil
}
