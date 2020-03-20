package base

// WithBody is something that has a configuration to allow insecure HTTP connections
type WithAllowInsecure interface {
	GetAllowInsecure() bool
}

// WithBody is something that has a body
type WithBody interface {
	GetBody() (string, error)
}

// WithHeaders is something that has headers
type WithHeaders interface {
	GetHeaders() map[string][]string
}

// WithMethod is something that has an HTTP method
type WithMethod interface {
	GetMethod() string
}

// WithValues is something that has values
type WithValues interface {
	GetValues() map[string][]string
}
