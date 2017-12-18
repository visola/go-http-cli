package daemon

import (
	"github.com/visola/go-http-cli/options"
	"github.com/visola/go-http-cli/request"
)

// ExecuteRequestResponse stores the request and response to be passed back when the daemon executes
// an HTTP request.
type ExecuteRequestResponse struct {
	RequestOptions *options.RequestOptions
	HTTPResponse   *request.HTTPResponse
}

// HandshakeResponse is the response sent by the daemon when someone is checking if it's up.
type HandshakeResponse struct {
	MajorVersion int8
	MinorVersion int8
}
