package profile

// Options that can come from a profile file.
type Options struct {
	AllowInsecure bool
	BaseURL       string
	Headers       map[string][]string
	NamedRequest  map[string]NamedRequest
	Variables     map[string]string
}

// GetAllowInsecure returns if this option allow insecure HTTP connections
func (ops Options) GetAllowInsecure() bool {
	return ops.AllowInsecure
}

// GetHeaders returns the headers set in this option
func (ops Options) GetHeaders() map[string][]string {
	return ops.Headers
}

// MergeOptions merges all options passed in into a final Options object.
func MergeOptions(profiles []Options) Options {
	baseURL := ""
	headers := make(map[string][]string)
	insecure := false
	requests := make(map[string]NamedRequest)
	variables := make(map[string]string)

	// Merge all profiles
	for _, profile := range profiles {
		if profile.BaseURL != "" {
			baseURL = profile.BaseURL
		}

		insecure = insecure || profile.AllowInsecure

		for header, values := range profile.Headers {
			headers[header] = append(headers[header], values...)
		}

		for variable, value := range profile.Variables {
			variables[variable] = value
		}

		for requestName, requestConfiguration := range profile.NamedRequest {
			requests[requestName] = requestConfiguration
		}
	}

	return Options{
		AllowInsecure: insecure,
		BaseURL:       baseURL,
		Headers:       headers,
		NamedRequest:  requests,
		Variables:     variables,
	}
}
