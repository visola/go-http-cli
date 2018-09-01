package request

// ExecutionOptions represent the options to be passed for the request executor.
type ExecutionOptions struct {
	FileToUpload    string
	FollowLocation  bool
	MaxRedirect     int
	PostProcessFile string
	ProfileNames    []string
	RequestName     string
	Request         Request
	Variables       map[string]string
}
