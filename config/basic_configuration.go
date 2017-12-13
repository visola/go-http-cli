package config

// BasicConfiguration is a simple implementation of the Configuration interface that can be used
// for un/marshalling purposes.
type BasicConfiguration struct {
	BaseURLField   string
	HeadersField   map[string][]string
	BodyField      string
	MethodField    string
	URLField       string
	VariablesField map[string]string
}

// BaseURL test implementation
func (conf BasicConfiguration) BaseURL() string {
	return conf.BaseURLField
}

// Headers test implementation
func (conf BasicConfiguration) Headers() map[string][]string {
	return conf.HeadersField
}

// Body test implementation
func (conf BasicConfiguration) Body() string {
	return conf.BodyField
}

// Method test implementation
func (conf BasicConfiguration) Method() string {
	return conf.MethodField
}

// URL test implementation
func (conf BasicConfiguration) URL() string {
	return conf.URLField
}

// Variables test implementation
func (conf BasicConfiguration) Variables() map[string]string {
	return conf.VariablesField
}

// ToBasicConfiguration converts any configuration to a basic configuration
func ToBasicConfiguration(configuration Configuration) *BasicConfiguration {
	return &BasicConfiguration{
		BaseURLField:   configuration.BaseURL(),
		BodyField:      configuration.Body(),
		HeadersField:   configuration.Headers(),
		MethodField:    configuration.Method(),
		URLField:       configuration.URL(),
		VariablesField: configuration.Variables(),
	}
}
