package daemon

import "github.com/visola/go-http-cli/pkg/request"

// HandshakeResponse is the response sent by the daemon when someone is checking if it's up.
type HandshakeResponse struct {
	MajorVersion int8
	MinorVersion int8
}

// RequestExecution is the response from the daemon when executing a request.
type RequestExecution struct {
	RequestResponses []request.ExecutedRequestResponse
	ErrorMessage     string
}
