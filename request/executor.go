package request

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
)

const defaultMaxRedirectCount = 10

// ExecuteRequest executes an HTTP request based on the specified options.
func ExecuteRequest(executionOptions ExecutionOptions) ([]ExecutedRequestResponse, error) {
	reqToExecute, getRequestError := getRequestFrom(executionOptions)

	maxRedirectCount := executionOptions.MaxRedirect
	if maxRedirectCount == 0 {
		maxRedirectCount = defaultMaxRedirectCount
	}

	if getRequestError != nil {
		return nil, getRequestError
	}

	requestsToExecute := make([]Request, 1)
	requestsToExecute[0] = reqToExecute

	redirectCount := 0

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	result := make([]ExecutedRequestResponse, 0)

	for {
		currentRequest := requestsToExecute[0]
		requestsToExecute = requestsToExecute[1:]

		httpRequest, currentConfiguredRequest, httpRequestErr := BuildRequest(currentRequest, executionOptions.ProfileNames, executionOptions.Variables)
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

		response := Response{
			StatusCode: httpResponse.StatusCode,
			Status:     httpResponse.Status,
			Headers:    headers,
			Body:       string(bodyBytes),
			Protocol:   fmt.Sprintf("%d.%d", httpResponse.ProtoMajor, httpResponse.ProtoMinor),
		}

		result = append(result, ExecutedRequestResponse{
			Request:  *currentConfiguredRequest,
			Response: response,
		})

		if shouldRedirect(response.StatusCode) && executionOptions.FollowLocation == true {
			redirectCount++

			if redirectCount > maxRedirectCount {
				return result, fmt.Errorf("Max number of redirects reached: %d", maxRedirectCount)
			}

			redirectRequest, redirectError := buildRedirect(*httpResponse)
			if redirectError != nil {
				return result, redirectError
			}
			requestsToExecute = append(requestsToExecute, *redirectRequest)
		}

		// If nothing else to execute, break
		if len(requestsToExecute) == 0 {
			break
		}
	}

	return result, nil
}

func buildRedirect(response http.Response) (*Request, error) {
	newLocation, responseError := response.Location()
	if responseError != nil {
		return nil, responseError
	}
	return &Request{
		URL: newLocation.String(),
	}, nil
}

func getRequestFrom(executionOptions ExecutionOptions) (Request, error) {
	var req Request
	if executionOptions.RequestName != "" {
		requestOptions, err := profile.LoadRequestOptions(executionOptions.RequestName, executionOptions.ProfileNames)

		if err != nil {
			return req, err
		}

		req = Request{
			Body:    requestOptions.Body,
			Headers: requestOptions.Headers,
			Method:  requestOptions.Method,
			URL:     requestOptions.URL,
		}

		req.Merge(executionOptions.Request)
	} else {
		req = executionOptions.Request
	}

	return req, nil
}

func shouldRedirect(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently ||
		statusCode == http.StatusFound ||
		statusCode == http.StatusSeeOther
}

func storeCookies(httpRequest http.Request, httpResponse http.Response) error {
	session, sessionErr := session.Get(*httpRequest.URL)

	if sessionErr != nil {
		return sessionErr
	}

	session.Jar.SetCookies(httpRequest.URL, httpResponse.Cookies())

	return nil
}
