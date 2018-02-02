package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	myioutil "github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
	mystrings "github.com/visola/go-http-cli/strings"
)

// BuildRequest builds a Request from a Configuration.
func BuildRequest(unconfiguredRequest Request, requestName string, executionOptions ExecutionOptions) (*http.Request, *Request, error) {
	profiles, profileError := loadProfiles(executionOptions.ProfileNames)

	if profileError != nil {
		return nil, nil, profileError
	}

	configuredRequest, configError := configureRequest(unconfiguredRequest, requestName, profiles, executionOptions)

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
		req.Body = myioutil.CreateCloseableBufferString(configuredRequest.Body)
	}

	return req, configuredRequest, nil
}

func configureRequest(unconfiguredRequest Request, requestName string, profiles []profile.Options, executionOptions ExecutionOptions) (*Request, error) {
	mergedProfile := profile.MergeOptions(profiles)
	mergeWithPassedIn(executionOptions, unconfiguredRequest, mergedProfile)

	var requestOptions profile.RequestOptions
	var exists bool
	if requestOptions, exists = mergedProfile.RequestOptions[requestName]; requestName != "" && !exists {
		return nil, fmt.Errorf("Request with name %s not found", requestName)
	}

	body, loadBodyErr := loadBody(unconfiguredRequest, requestOptions, executionOptions)
	if loadBodyErr != nil {
		return nil, loadBodyErr
	}

	unconfiguredRequest = mergeRequests(unconfiguredRequest, requestOptions)

	method := unconfiguredRequest.Method
	if method == "" {
		if body == "" {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
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
		Body:    replaceVariables(body, mergedProfile.Variables),
		Cookies: session.Jar.Cookies(url),
		Headers: replaceVariablesInHeaders(mergedProfile.Headers, mergedProfile.Variables),
		Method:  method,
		URL:     urlString,
	}, nil
}

func loadBody(unconfiguredRequest Request, requestOptions profile.RequestOptions, executionOptions ExecutionOptions) (string, error) {
	// Passed in body overrides everything
	if unconfiguredRequest.Body != "" {
		return unconfiguredRequest.Body, nil
	}

	// Passed in file to upload overrides everything
	if executionOptions.FileToUpload != "" {
		data, loadError := ioutil.ReadFile(executionOptions.FileToUpload)
		if loadError != nil {
			return "", loadError
		}

		return string(data), nil
	}

	if requestOptions.Body != "" {
		return requestOptions.Body, nil
	}

	// If there's a file to upload from profile, load it
	if requestOptions.FileToUpload != "" {
		path := requestOptions.FileToUpload
		// If not absolute, it's relative to the profiles dir
		if !filepath.IsAbs(path) {
			profileDir, profileDirError := profile.GetProfilesDir()
			if profileDirError != nil {
				return "", profileDirError
			}

			path = filepath.Join(profileDir, requestOptions.FileToUpload)
		}

		data, loadErr := ioutil.ReadFile(path)
		if loadErr != nil {
			return "", loadErr
		}

		return string(data), nil
	}

	return "", nil
}

func loadProfiles(profileNames []string) ([]profile.Options, error) {
	profiles := make([]profile.Options, len(profileNames))

	for index, profileName := range profileNames {
		profile, profileError := profile.LoadProfile(profileName)
		if profileError != nil {
			return nil, profileError
		}
		profiles[index] = *profile
	}

	return profiles, nil
}

func mergeRequests(unconfiguredRequest Request, requestOptions profile.RequestOptions) Request {
	loadedRequest := Request{
		Body:    requestOptions.Body,
		Headers: requestOptions.Headers,
		Method:  requestOptions.Method,
		URL:     requestOptions.URL,
	}

	newRequest := Request{}
	newRequest.Merge(loadedRequest)
	newRequest.Merge(unconfiguredRequest)
	return newRequest
}

func mergeWithPassedIn(executionOptions ExecutionOptions, unconfiguredRequest Request, profile profile.Options) {
	// Merge the passed in variables
	for variable, value := range executionOptions.Variables {
		profile.Variables[variable] = value
	}

	// Merge all headers
	for header, values := range unconfiguredRequest.Headers {
		profile.Headers[header] = append(profile.Headers[header], values...)
	}
}

func replaceVariablesInHeaders(headers map[string][]string, variables map[string]string) map[string][]string {
	result := make(map[string][]string)
	for header, values := range headers {
		newValues := make([]string, len(values))
		for index, value := range values {
			newValues[index] = mystrings.ParseExpression(value, variables)
		}
		result[header] = newValues
	}
	return result
}

func replaceVariables(value string, variables map[string]string) string {
	return mystrings.ParseExpression(value, variables)
}
