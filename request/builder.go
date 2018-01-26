package request

import (
	"net/http"
	"net/url"

	"github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
	"github.com/visola/go-http-cli/strings"
)

// BuildRequest builds a Request from a Configuration.
func BuildRequest(unconfiguredRequest Request, profileNames []string, variables map[string]string) (*http.Request, *Request, error) {
	profiles, profileError := loadProfiles(profileNames)

	if profileError != nil {
		return nil, nil, profileError
	}

	configuredRequest, configError := configureRequest(unconfiguredRequest, profiles, variables)

	if configError != nil {
		return nil, nil, configError
	}

	req, reqErr := http.NewRequest(configuredRequest.Method, configuredRequest.URL, nil)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	for _, cookie := range configuredRequest.Cookies {
		req.AddCookie(cookie)
	}

	for k, vs := range configuredRequest.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if configuredRequest.Body != "" {
		req.Body = ioutil.CreateCloseableBufferString(configuredRequest.Body)
	}

	return req, configuredRequest, nil
}

func configureRequest(unconfiguredRequest Request, profiles []profile.Options, passedVariables map[string]string) (*Request, error) {
	mergedProfile := profile.MergeOptions(profiles)

	method := unconfiguredRequest.Method

	if method == "" {
		if unconfiguredRequest.Body == "" {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
	}

	// Merge the passed in variables
	for variable, value := range passedVariables {
		mergedProfile.Variables[variable] = value
	}

	// Merge all headers
	for header, values := range unconfiguredRequest.Headers {
		mergedProfile.Headers[header] = append(mergedProfile.Headers[header], values...)
	}

	urlString := replaceVariables(ParseURL(mergedProfile.BaseURL, unconfiguredRequest.URL), mergedProfile.Variables)

	url, urlError := url.Parse(urlString)

	if urlError != nil {
		return nil, urlError
	}

	session, sessionErr := session.Get(*url)

	if sessionErr != nil {
		return nil, sessionErr
	}

	return &Request{
		Body:    replaceVariables(unconfiguredRequest.Body, mergedProfile.Variables),
		Cookies: session.Jar.Cookies(url),
		Headers: replaceVariablesInHeaders(mergedProfile.Headers, mergedProfile.Variables),
		Method:  method,
		URL:     urlString,
	}, nil
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

func replaceVariablesInHeaders(headers map[string][]string, variables map[string]string) map[string][]string {
	result := make(map[string][]string)
	for header, values := range headers {
		newValues := make([]string, len(values))
		for index, value := range values {
			newValues[index] = strings.ParseExpression(value, variables)
		}
		result[header] = newValues
	}
	return result
}

func replaceVariables(value string, variables map[string]string) string {
	return strings.ParseExpression(value, variables)
}
