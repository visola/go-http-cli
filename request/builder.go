package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	myioutil "github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	mystrings "github.com/visola/go-http-cli/strings"
)

const jsonMimeType = "application/json"
const urlEncodedMimeType = "application/x-www-form-urlencoded"

var bodyBuilderContentTypes = [...]string{
	urlEncodedMimeType,
	jsonMimeType,
}

// BuildRequest builds an http.Request from a configured request.Request
func BuildRequest(configuredRequest Request) (*http.Request, error) {
	req, reqErr := http.NewRequest(configuredRequest.Method, configuredRequest.URL, nil)
	if reqErr != nil {
		return nil, reqErr
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

	return req, nil
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
		Values:  requestOptions.Values,
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

func setMethod(initialMethod string, headers map[string][]string, hasValues bool, body string) string {
	method := initialMethod

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
		if hasValues {
			contenType := getContentType(headers)
			for _, knownType := range bodyBuilderContentTypes {
				if strings.HasPrefix(contenType, knownType) {
					return http.MethodPost
				}
			}
		}
	}

	return method
}
