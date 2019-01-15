package request

// ExecutionContext represent the options to be passed for the request executor.
type ExecutionContext struct {
	FollowLocation  bool
	MaxRedirect     int
	PostProcessCode PostProcessSourceCode
	ProfileNames    []string
	Request         Request
	Variables       map[string]string
}
