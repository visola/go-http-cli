package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	myioutil "github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
	mystrings "github.com/visola/go-http-cli/strings"
)

const jsonMimeType = "application/json"
const urlEncodedMimeType = "application/x-www-form-urlencoded"

var bodyBuilderContentTypes = [...]string{
	urlEncodedMimeType,
	jsonMimeType,
}

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
	mergeVariablesIn(mergedProfile, executionOptions)
	mergeHeadersIn(mergedProfile, unconfiguredRequest.Headers)

	var requestOptions profile.RequestOptions
	var exists bool
	if requestOptions, exists = mergedProfile.RequestOptions[requestName]; requestName != "" && !exists {
		return nil, fmt.Errorf("Request with name %s not found", requestName)
	}

	mergeHeadersIn(mergedProfile, requestOptions.Headers)

	body, loadBodyErr := loadBody(unconfiguredRequest, requestOptions, executionOptions)
	if loadBodyErr != nil {
		return nil, loadBodyErr
	}

	unconfiguredRequest = mergeRequests(unconfiguredRequest, requestOptions)
	method := setMethod(unconfiguredRequest, body)

	urlString := replaceVariables(ParseURL(mergedProfile.BaseURL, unconfiguredRequest.URL), mergedProfile.Variables)

	parsedURL, urlError := url.Parse(urlString)
	if urlError != nil {
		return nil, urlError
	}

	if len(unconfiguredRequest.Values) > 0 {
		replacedValues := replaceVariablesInMapOfArrayOfStrings(unconfiguredRequest.Values, mergedProfile.Variables)
		if method == http.MethodGet {
			parsedURL.RawQuery = encodeValues(replacedValues)
		} else {
			contentType := getContentType(mergedProfile.Headers)
			body = createBody(contentType, replacedValues)
			if contentType == "" {
				mergedProfile.Headers["Content-Type"] = []string{jsonMimeType}
			}
		}
	}

	session, sessionErr := session.Get(*parsedURL)
	if sessionErr != nil {
		return nil, sessionErr
	}

	return &Request{
		Body:    replaceVariables(body, mergedProfile.Variables),
		Cookies: session.Jar.Cookies(parsedURL),
		Headers: replaceVariablesInMapOfArrayOfStrings(mergedProfile.Headers, mergedProfile.Variables),
		Method:  method,
		URL:     parsedURL.String(),
	}, nil
}

func buildJSON(values map[string][]string) string {
	toEncode := make(map[string]string)
	for key, valuesForKey := range values {
		for _, value := range valuesForKey {
			toEncode[key] = value
		}
	}

	// Ignore this error, encoding map to JSON should never fail
	jsonBytes, _ := json.Marshal(toEncode)
	return string(jsonBytes)
}

func createBody(contentType string, values map[string][]string) string {
	if contentType == "" || strings.HasSuffix(strings.TrimSpace(contentType), jsonMimeType) {
		return buildJSON(values)
	} else if strings.HasSuffix(strings.TrimSpace(contentType), urlEncodedMimeType) {
		return encodeValues(values)
	}

	return fmt.Sprintf("Unsupported body type: %s", contentType)
}

func encodeValues(values map[string][]string) string {
	vals := url.Values{}
	for name, valuesForKey := range values {
		for _, value := range valuesForKey {
			vals.Add(name, value)
		}
	}
	return vals.Encode()
}

func getContentType(headers map[string][]string) string {
	for name, values := range headers {
		if strings.ToLower(strings.TrimSpace(name)) == "content-type" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
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
		profiles[index] = profile
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

func mergeVariablesIn(profile profile.Options, executionOptions ExecutionOptions) {
	for variable, value := range executionOptions.Variables {
		profile.Variables[variable] = value
	}
}

func mergeHeadersIn(profile profile.Options, headers map[string][]string) {
	for header, values := range headers {
		profile.Headers[header] = values
	}
}

func replaceVariablesInMapOfArrayOfStrings(headers map[string][]string, variables map[string]string) map[string][]string {
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

func setMethod(req Request, body string) string {
	method := req.Method

	if method == "" {
		if body == "" {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
	}

	// If still empty
	if method == "" || method == http.MethodGet {
		// If there are values, check if they should go in the body
		if len(req.Values) != 0 {
			contenType := getContentType(req.Headers)
			for _, knownType := range bodyBuilderContentTypes {
				if strings.HasPrefix(contenType, knownType) {
					return http.MethodPost
				}
			}
		}
	}

	return method
}
