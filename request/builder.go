package request

import (
	"net/http"

	"github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
)

// BuildRequest builds a Request from a Configuration.
func BuildRequest(unconfiguredRequest Request, profileNames []string, variables map[string]string) (*http.Request, *Request, error) {
	profiles, profileError := loadProfiles(profileNames)

	if profileError != nil {
		return nil, nil, profileError
	}

	configuredRequest := configureRequest(unconfiguredRequest, profiles, variables)

	req, reqErr := http.NewRequest(configuredRequest.Method, configuredRequest.URL, nil)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	for k, vs := range configuredRequest.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if configuredRequest.Body != "" {
		req.Body = ioutil.CreateCloseableBufferString(configuredRequest.Body)
	}

	return req, &configuredRequest, nil
}

func configureRequest(unconfiguredRequest Request, profiles []profile.Options, passedVariables map[string]string) Request {
	mergedProfile := profile.MergeOptions(profiles)

	// Merge the passed in variables
	for variable, value := range passedVariables {
		mergedProfile.Variables[variable] = value
	}

	// Merge all headers
	for header, values := range unconfiguredRequest.Headers {
		mergedProfile.Headers[header] = append(mergedProfile.Headers[header], values...)
	}

	url := ParseURL(mergedProfile.BaseURL, unconfiguredRequest.URL, mergedProfile.Variables)

	return Request{
		Body:    unconfiguredRequest.Body,
		Headers: mergedProfile.Headers,
		Method:  unconfiguredRequest.Method,
		URL:     url,
	}
}

func loadProfiles(profileNames []string) ([]profile.Options, error) {
	profiles := make([]profile.Options, len(profileNames))

	for _, profileName := range profileNames {
		profile, profileError := profile.LoadProfile(profileName)
		if profileError != nil {
			return nil, profileError
		}
		profiles = append(profiles, *profile)
	}

	return profiles, nil
}
