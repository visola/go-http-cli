package options

// RequestOptions stores data required to configure a request to be executed
type RequestOptions struct {
	Body      string
	Headers   map[string][]string
	Method    string
	Profiles  []string
	URL       string
	Variables map[string]string
}
