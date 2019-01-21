package request

import "github.com/visola/go-http-cli/pkg/session"

// ExecutionContext represent the options to be passed for the request executor.
type ExecutionContext struct {
	FollowLocation  bool
	MaxRedirect     int
	PostProcessCode PostProcessSourceCode
	ProfileNames    []string
	Request         Request
	Session         *session.Session
	Variables       map[string]string
}
