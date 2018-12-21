package request

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/visola/go-http-cli/session"
)

const defaultMaxRedirectCount = 10

// ExecuteRequestLoop executes HTTP requests based on the passed in options until there're no more
// requests to be executed.
func ExecuteRequestLoop(executionOptions ExecutionOptions) ([]ExecutedRequestResponse, error) {
	maxRedirectCount := executionOptions.MaxRedirect
	if maxRedirectCount == 0 {
		maxRedirectCount = defaultMaxRedirectCount
	}

	requestsToExecute := make([]Request, 1)
	requestsToExecute[0] = executionOptions.Request

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	result := make([]ExecutedRequestResponse, 0)
	redirectCount := 0
	for {
		currentConfiguredRequest := requestsToExecute[0]
		requestsToExecute = requestsToExecute[1:]

		currentConfiguredRequest, processError := replaceRequestVariables(currentConfiguredRequest, executionOptions.ProfileNames, executionOptions.Variables)
		if processError != nil {
			return nil, processError
		}

		response, executeErr := executeRequest(client, currentConfiguredRequest)

		result = append(result, ExecutedRequestResponse{
			Request:  currentConfiguredRequest,
			Response: *response,
		})

		if executeErr != nil {
			return result, executeErr
		}

		if shouldRedirect(response.StatusCode) && executionOptions.FollowLocation == true {
			redirectCount++

			if redirectCount > maxRedirectCount {
				return result, fmt.Errorf("Max number of redirects reached: %d", maxRedirectCount)
			}

			redirectRequest := buildRedirect(response)
			requestsToExecute = append(requestsToExecute, *redirectRequest)
		}

		// If nothing else to execute, break
		if len(requestsToExecute) == 0 {
			break
		}
	}

	return result, nil
}

func buildRedirect(response *Response) *Request {
	location := response.Headers["Location"]
	if len(location) > 0 && location[0] != "" {
		return &Request{
			URL: location[0],
		}
	}

	return nil
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
