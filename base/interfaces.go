package base

// WithBody is something that has a body
type WithBody interface {
	GetBody() (string, error)
}

// WithHeaders is something that has headers
type WithHeaders interface {
	GetHeaders() map[string][]string
}

// WithValues is something that has values
type WithValues interface {
	GetValues() map[string][]string
}
