package request

// Request stores data required to configure a request to be executed
type Request struct {
	Body    string
	Headers map[string][]string
	Method  string
	URL     string
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
