package request

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/visola/go-http-cli/config"
	"github.com/visola/go-http-cli/options"
)

// HTTPResponse is the response from the daemon after executing a request
type HTTPResponse struct {
	Body       string
	Headers    map[string][]string
	Protocol   string
	StatusCode int
	Status     string
}

// ExecuteRequest loads all the profile information and other related data associated with the
// passed in options and execute an HTTP request based on the parsed options.
func ExecuteRequest(options options.RequestOptions) (*HTTPResponse, error) {
	configuration, configError := config.Parse(options)

	if configError != nil {
		return nil, configError
	}

	request, requestErr := BuildRequest(configuration)
	if requestErr != nil {
		return nil, requestErr
	}

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, responseErr
	}

	bodyBytes, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return nil, readErr
	}

	body := ""
	if len(bodyBytes) != 0 {
		body = string(bodyBytes)
	}

	headers := make(map[string][]string)
	for k, vs := range response.Header {
		headers[k] = append(headers[k], vs...)
	}

	return &HTTPResponse{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		Headers:    headers,
		Body:       body,
		Protocol:   fmt.Sprintf("%d.%d", response.ProtoMajor, response.ProtoMinor),
	}, nil
}
