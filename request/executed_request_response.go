package request

// ExecutedRequestResponse represents a pair of request and the response that was returned from its
// execution.
type ExecutedRequestResponse struct {
	Request  Request
	Response Response
}
