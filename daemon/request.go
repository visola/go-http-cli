package daemon

import (
	"github.com/visola/go-http-cli/request"
)

// Request represents the data that needs to be passed to the daemon in order to execute an HTTP request.
type Request struct {
	Body      string
	Headers   map[string][]string
	Method    string
	Profiles  []string
	URL       string
	Variables map[string]string
}

// ToRequest transform a daemon.Request into a request.Request
func (req Request) ToRequest() request.Request {
	return request.Request{
		Body:    req.Body,
		Headers: req.Headers,
		Method:  req.Method,
		URL:     req.URL,
	}
}
