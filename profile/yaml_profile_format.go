package profile

// Used to unmarshal data from YAML files
type yamlProfileFormat struct {
	BaseURL   string `yaml:"baseURL"`
	Headers   map[string]arrayOrString
	Variables map[string]string
	Requests  []requestConfiguration
}

// Used to unmarshal request options from yaml files
type requestConfiguration struct {
	Body    string
	Headers map[string]arrayOrString
	Method  string
	Name    string
	URL     string
}

func (loadedProfile yamlProfileFormat) toOptions() *Options {
	return &Options{
		BaseURL:        loadedProfile.BaseURL,
		Headers:        toMapOfArrayOfStrings(loadedProfile.Headers),
		RequestOptions: toMapOfRequestOptions(loadedProfile.Requests),
		Variables:      loadedProfile.Variables,
	}
}

func toMapOfRequestOptions(requestConfigurations []requestConfiguration) map[string]RequestOptions {
	result := make(map[string]RequestOptions)
	for _, requestConfiguration := range requestConfigurations {
		result[requestConfiguration.Name] = RequestOptions{
			Body:    requestConfiguration.Body,
			Headers: toMapOfArrayOfStrings(requestConfiguration.Headers),
			Method:  requestConfiguration.Method,
			URL:     requestConfiguration.URL,
		}
	}

	return result
}
