package profile

// Used to unmarshal data from YAML files
type yamlProfileFormat struct {
	BaseURL   string `yaml:"baseURL"`
	Headers   map[string]arrayOrString
	Variables map[string]string
	Requests  map[string]requestConfiguration
}

// Used to unmarshal request options from yaml files
type requestConfiguration struct {
	Body    string
	Headers map[string]arrayOrString
	Method  string
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

func toMapOfRequestOptions(requestConfigurations map[string]requestConfiguration) map[string]RequestOptions {
	result := make(map[string]RequestOptions)
	for name, requestConfiguration := range requestConfigurations {
		result[name] = RequestOptions{
			Body:    requestConfiguration.Body,
			Headers: toMapOfArrayOfStrings(requestConfiguration.Headers),
			Method:  requestConfiguration.Method,
			URL:     requestConfiguration.URL,
		}
	}

	return result
}
