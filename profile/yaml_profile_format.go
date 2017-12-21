package profile

// Used to unmarshal data from YAML files
type yamlProfileFormat struct {
	BaseURL   string `yaml:"baseURL"`
	Headers   map[string]arrayOrString
	Variables map[string]string
}

func (loadedProfile yamlProfileFormat) toOptions() *Options {
	return &Options{
		BaseURL:   loadedProfile.BaseURL,
		Headers:   toMapOfArrayOfStrings(loadedProfile.Headers),
		Variables: loadedProfile.Variables,
	}
}
