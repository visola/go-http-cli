package config

type hierarchicalConfigurationFormat struct {
	configurations []Configuration
}

func (conf hierarchicalConfigurationFormat) BaseURL() string {
	var result string
	for _, subConfig := range conf.configurations {
		if subConfig.BaseURL() != "" {
			result = subConfig.BaseURL()
		}
	}
	return result
}

func (conf hierarchicalConfigurationFormat) Headers() map[string][]string {
	result := make(map[string][]string)
	for _, subConfig := range conf.configurations {
		for k, v := range subConfig.Headers() {
			result[k] = v
		}
	}
	return result
}

func (conf hierarchicalConfigurationFormat) Body() string {
	var result string
	for _, subConfig := range conf.configurations {
		if subConfig.Body() != "" {
			result = subConfig.Body()
		}
	}
	return result
}

func (conf hierarchicalConfigurationFormat) Method() string {
	var result string
	for _, subConfig := range conf.configurations {
		if subConfig.Method() != "" {
			result = subConfig.Method()
		}
	}
	return result
}

func (conf hierarchicalConfigurationFormat) URL() string {
	var result string
	for _, subConfig := range conf.configurations {
		if subConfig.URL() != "" {
			result = subConfig.URL()
		}
	}
	return result
}

func (conf hierarchicalConfigurationFormat) Variables() map[string]string {
	result := make(map[string]string)
	for _, subConfig := range conf.configurations {
		for k, v := range subConfig.Variables() {
			result[k] = v
		}
	}
	return result
}
