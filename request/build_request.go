package request

import (
	"net/http"
	"net/url"

	"github.com/visola/go-http-cli/ioutil"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/variables"
)

// BuildRequest builds an http.Request from a configured request.Request
func BuildRequest(configuredRequest Request, profileNames []string, passedInVariables map[string]string) (*http.Request, error) {
	processedRequest, processError := replaceRequestVariables(configuredRequest, profileNames, passedInVariables)
	if processError != nil {
		return nil, processError
	}

	parsedURL, urlError := url.Parse(processedRequest.URL)
	if urlError != nil {
		return nil, urlError
	}

	parsedURL.RawQuery = encodeValues(processedRequest.QueryParams)

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

	req.Body = ioutil.CreateCloseableBufferString(processedRequest.Body)
	return req, nil
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
	configuredRequest.QueryParams = replaceVariablesInMapOfArrayOfStrings(configuredRequest.QueryParams, finalVariableSet)

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
