package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ExecuteRequest executes an HTTP request based on the specified options.
func ExecuteRequest(request Request, profileNames []string, variables map[string]string) ([]ExecutedRequestResponse, error) {
	toExecute := make([]Request, 1)
	toExecute[0] = request

	client := &http.Client{
		// Do not auto-follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	result := make([]ExecutedRequestResponse, 0)

	for {
		executing := toExecute[0]
		toExecute = toExecute[1:]

		httpRequest, configuredRequest, httpRequestErr := BuildRequest(executing, profileNames, variables)
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

		body := ""
		if len(bodyBytes) != 0 {
			body = string(bodyBytes)
		}

		headers := make(map[string][]string)
		for k, vs := range httpResponse.Header {
			headers[k] = append(headers[k], vs...)
		}

		response := Response{
			StatusCode: httpResponse.StatusCode,
			Status:     httpResponse.Status,
			Headers:    headers,
			Body:       body,
			Protocol:   fmt.Sprintf("%d.%d", httpResponse.ProtoMajor, httpResponse.ProtoMinor),
		}

		result = append(result, ExecutedRequestResponse{
			Request:  *configuredRequest,
			Response: response,
		})

		if response.StatusCode == http.StatusMovedPermanently || response.StatusCode == http.StatusFound || response.StatusCode == http.StatusSeeOther {
			redirectRequest, redirectError := buildRedirect(*httpResponse)
			if redirectError != nil {
				return result, redirectError
			}
			toExecute = append(toExecute, *redirectRequest)
		}

		// If nothing else to execute, break
		if len(toExecute) == 0 {
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
