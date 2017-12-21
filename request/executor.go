package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ExecuteRequest executes an HTTP request based on the specified options.
func ExecuteRequest(request Request, profileNames []string, variables map[string]string) (*ExecutedRequestResponse, error) {
	httpRequest, configuredRequest, httpRequestErr := BuildRequest(request, profileNames, variables)
	if httpRequestErr != nil {
		return nil, httpRequestErr
	}

	client := &http.Client{}
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

	return &ExecutedRequestResponse{
		Request:  *configuredRequest,
		Response: response,
	}, nil
}
