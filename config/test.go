package config

// TestConfiguration is a simple implementation of the Configuration interface that can be used
// for testing purposes.
type TestConfiguration struct {
	TestBaseURL   string
	TestHeaders   map[string][]string
	TestBody      string
	TestMethod    string
	TestURL       string
	TestVariables map[string]string
}

// BaseURL test implementation
func (conf TestConfiguration) BaseURL() string {
	return conf.TestBaseURL
}

// Headers test implementation
func (conf TestConfiguration) Headers() map[string][]string {
	return conf.TestHeaders
}

// Body test implementation
func (conf TestConfiguration) Body() string {
	return conf.TestBody
}

// Method test implementation
func (conf TestConfiguration) Method() string {
	return conf.TestMethod
}

// URL test implementation
func (conf TestConfiguration) URL() string {
	return conf.TestURL
}

// Variables test implementation
func (conf TestConfiguration) Variables() map[string]string {
	return conf.TestVariables
}
