package daemon

// ExecuteRequestResponse is the response from the daemon after executing a request
type ExecuteRequestResponse struct {
	Body       string
	Headers    map[string][]string
	Protocol   string
	StatusCode int
	Status     string
}

// HandshakeResponse is the response sent by the daemon when someone is checking if it's up.
type HandshakeResponse struct {
	MajorVersion int8
	MinorVersion int8
}
