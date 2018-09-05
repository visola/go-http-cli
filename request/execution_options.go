package request

// ExecutionOptions represent the options to be passed for the request executor.
type ExecutionOptions struct {
	FollowLocation  bool
	MaxRedirect     int
	PostProcessFile string
	ProfileNames    []string
	Request         Request
	Variables       map[string]string
}
