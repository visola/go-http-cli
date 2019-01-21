package request

// ExecutedRequestResponse represents a pair of request and the response that was returned from its
// execution. It also includes any output and/or error generated during post processing.
type ExecutedRequestResponse struct {
	Request           Request
	Response          Response
	PostProcessError  string
	PostProcessOutput string
}
