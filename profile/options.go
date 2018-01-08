package profile

// Options that can come from a profile file.
type Options struct {
	BaseURL        string
	Headers        map[string][]string
	RequestOptions map[string]RequestOptions
	Variables      map[string]string
}

// MergeOptions merges all options passed in into a final Options object.
func MergeOptions(profiles []Options) Options {
	baseURL := ""
	headers := make(map[string][]string)
	requests := make(map[string]RequestOptions)
	variables := make(map[string]string)

	// Merge all profiles
	for _, profile := range profiles {
		if profile.BaseURL != "" {
			baseURL = profile.BaseURL
		}

		for header, values := range profile.Headers {
			headers[header] = append(headers[header], values...)
		}

		for variable, value := range profile.Variables {
			variables[variable] = value
		}

		for requestName, requestConfiguration := range profile.RequestOptions {
			requests[requestName] = requestConfiguration
		}
	}

	return Options{
		BaseURL:        baseURL,
		Headers:        headers,
		RequestOptions: requests,
		Variables:      variables,
	}
}
