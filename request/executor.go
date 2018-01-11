package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const maxRedirect = 10

// ExecuteRequest executes an HTTP request based on the specified options.
func ExecuteRequest(request Request, profileNames []string, variables map[string]string) ([]ExecutedRequestResponse, error) {
	requestsToExecute := make([]Request, 1)
	requestsToExecute[0] = request
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

		httpRequest, currentConfiguredRequest, httpRequestErr := BuildRequest(currentRequest, profileNames, variables)
		if httpRequestErr != nil {
			return nil, httpRequestErr
		}

		httpResponse, httpResponseErr := client.Do(httpRequest)
		if httpResponseErr != nil {
			return nil, httpResponseErr
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

		if shouldRedirect(response.StatusCode) {
			redirectCount++

			if redirectCount > maxRedirect {
				return result, fmt.Errorf("Max number of redirects reached: %d", maxRedirect)
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

func shouldRedirect(statusCode int) bool {
	return statusCode == http.StatusMovedPermanently ||
		statusCode == http.StatusFound ||
		statusCode == http.StatusSeeOther
}
