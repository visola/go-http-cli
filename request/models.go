package request

// Request stores data required to configure a request to be executed
type Request struct {
	Body    string
	Headers map[string][]string
	Method  string
	URL     string
}

// Merge merges information from another request into the original request, overwriting any data
// that is provided in toMerge.
func (original *Request) Merge(toMerge Request) {
	if toMerge.Body != "" {
		original.Body = toMerge.Body
	}

	for header, values := range toMerge.Headers {
		original.Headers[header] = values
	}

	if toMerge.Method != "" {
		original.Method = toMerge.Method
	}

	if toMerge.URL != "" {
		original.URL = toMerge.URL
	}
}

// Response is the response from the daemon after executing a request
type Response struct {
	Body       string
	Headers    map[string][]string
	Protocol   string
	StatusCode int
	Status     string
}

// ExecutedRequestResponse represents a pair of request and the response that was returned from its
// execution.
type ExecutedRequestResponse struct {
	Request  Request
	Response Response
}