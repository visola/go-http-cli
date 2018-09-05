package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	myioutil "github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/variables"
)

const jsonMimeType = "application/json"
const urlEncodedMimeType = "application/x-www-form-urlencoded"

var bodyBuilderContentTypes = [...]string{
	urlEncodedMimeType,
	jsonMimeType,
}

// BuildRequest builds an http.Request from a configured request.Request
func BuildRequest(configuredRequest Request, profileNames []string, passedInVariables map[string]string) (*http.Request, error) {
	processedRequest, processError := replaceRequestVariables(configuredRequest, profileNames, passedInVariables)
	if processError != nil {
		return nil, processError
	}

	fmt.Printf("Method before: %s\n", processedRequest.Method)
	processedRequest.Method = getMethod(processedRequest)
	fmt.Printf("Method after: %s\n", processedRequest.Method)

	parsedURL, urlError := url.Parse(processedRequest.URL)
	if urlError != nil {
		return nil, urlError
	}

	if processedRequest.Method == http.MethodGet {
		parsedURL.RawQuery = encodeValues(processedRequest.Values)
	}

	req, reqErr := http.NewRequest(processedRequest.Method, parsedURL.String(), nil)
	if reqErr != nil {
		return nil, reqErr
	}

	for _, cookie := range processedRequest.Cookies {
		req.AddCookie(cookie)
	}

	for k, vs := range processedRequest.Headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	req.Body = getBody(processedRequest)

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

func createBody(processedRequest Request) string {
	contentType := getContentType(processedRequest.Headers)
	if contentType == "" || strings.HasSuffix(strings.TrimSpace(contentType), jsonMimeType) {
		return buildJSON(processedRequest.Values)
	} else if strings.HasSuffix(strings.TrimSpace(contentType), urlEncodedMimeType) {
		return encodeValues(processedRequest.Values)
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

func getBody(processedRequest Request) *myioutil.CloseableByteBuffer {
	if processedRequest.Method == http.MethodGet {
		return myioutil.CreateCloseableBufferString("")
	}

	if processedRequest.Body != "" {
		return myioutil.CreateCloseableBufferString(processedRequest.Body)
	}

	if len(processedRequest.Values) > 0 {
		return myioutil.CreateCloseableBufferString(createBody(processedRequest))
	}

	return myioutil.CreateCloseableBufferString("")
}

func getContentType(headers map[string][]string) string {
	for name, values := range headers {
		if strings.ToLower(strings.TrimSpace(name)) == "content-type" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}

func getMethod(configuredRequest Request) string {
	method := configuredRequest.Method

	if method == "" {
		if configuredRequest.Body == "" {
			method = http.MethodGet
		} else {
			method = http.MethodPost
		}
	}

	// If still empty
	if method == "" || method == http.MethodGet {
		// If there are values, check if they should go in the body
		if len(configuredRequest.Values) > 0 {
			contenType := getContentType(configuredRequest.Headers)
			for _, knownType := range bodyBuilderContentTypes {
				if strings.HasPrefix(contenType, knownType) {
					return http.MethodPost
				}
			}
		}
	}

	return method
}

func replaceRequestVariables(configuredRequest Request, profileNames []string, passedInVariables map[string]string) (Request, error) {
	mergedProfiles, profileError := profile.LoadAndMergeProfiles(profileNames)
	if profileError != nil {
		return configuredRequest, profileError
	}

	finalVariableSet := mergeVariables(passedInVariables, mergedProfiles.Variables)

	configuredRequest.Body = variables.ReplaceVariables(configuredRequest.Body, finalVariableSet)
	configuredRequest.Headers = replaceVariablesInMapOfArrayOfStrings(configuredRequest.Headers, finalVariableSet)
	configuredRequest.URL = variables.ReplaceVariables(configuredRequest.URL, finalVariableSet)
	configuredRequest.Values = replaceVariablesInMapOfArrayOfStrings(configuredRequest.Values, finalVariableSet)

	return configuredRequest, nil
}

func mergeVariables(allVariables ...map[string]string) map[string]string {
	result := make(map[string]string)
	for i := len(allVariables) - 1; i >= 0; i-- {
		vars := allVariables[i]
		for key, val := range vars {
			result[key] = val
		}
	}
	return result
}

func replaceVariablesInMapOfArrayOfStrings(headers map[string][]string, context map[string]string) map[string][]string {
	result := make(map[string][]string)
	for header, values := range headers {
		newValues := make([]string, len(values))
		for index, value := range values {
			newValues[index] = variables.ReplaceVariables(value, context)
		}
		result[header] = newValues
	}
	return result
}
