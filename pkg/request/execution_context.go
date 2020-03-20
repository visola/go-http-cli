package request

import "github.com/visola/go-http-cli/pkg/session"

// ExecutionContext represent the options to be passed for the request executor.
type ExecutionContext struct {
	AllowInsecure    bool
	FollowLocation   bool
	MaxAddedRequests int
	MaxRedirect      int
	ProfileNames     []string
	Request          Request
	Session          *session.Session
	Variables        map[string]string
}
