package request

// ConfigureRequestOptions set of options that can be used to configure a request
type ConfigureRequestOptions struct {
	ProfileNames []string
	RequestName  string
	Values       map[string][]string
}

// GetValues returns the values from the configuration request
func (config *ConfigureRequestOptions) GetValues() map[string][]string {
	return config.Values
}

// ConfigureRequestOption an option that can be applied on a ConfigureRequestOptions
type ConfigureRequestOption func(configuration *ConfigureRequestOptions)

// AddProfiles adds a set of profiles to the configuration request
func AddProfiles(profileNames ...string) ConfigureRequestOption {
	return func(configuration *ConfigureRequestOptions) {
		configuration.ProfileNames = append(configuration.ProfileNames, profileNames...)
	}
}

// AddValues adds a set of values in the configuration request
func AddValues(values map[string][]string) ConfigureRequestOption {
	return func(configuration *ConfigureRequestOptions) {
		if configuration.Values == nil {
			configuration.Values = make(map[string][]string)
		}

		for valName, val := range values {
			configuration.Values[valName] = val
		}
	}
}

// SetRequestName sets the named request in the configuration request
func SetRequestName(requestName string) ConfigureRequestOption {
	return func(configuration *ConfigureRequestOptions) {
		configuration.RequestName = requestName
	}
}
