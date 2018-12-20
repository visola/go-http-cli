package request

import (
	"net/url"

	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/session"
	"github.com/visola/variables/variables"
)

func replaceRequestVariables(configuredRequest Request, profileNames []string, passedInVariables map[string]string) (Request, error) {
	mergedProfiles, profileError := profile.LoadAndMergeProfiles(profileNames)
	if profileError != nil {
		return configuredRequest, profileError
	}

	// Chicken-egg problem here, URL can have variables, but session can only be determined by URL
	finalVariableSet := mergeVariables(passedInVariables, mergedProfiles.Variables)
	configuredRequest.URL = variables.ReplaceVariables(configuredRequest.URL, finalVariableSet)

	parsedURL, parseURLErr := url.Parse(configuredRequest.URL)
	if parseURLErr != nil {
		return configuredRequest, parseURLErr
	}

	session, sessionError := session.Get(parsedURL.Hostname())
	if sessionError != nil {
		return configuredRequest, sessionError
	}

	finalVariableSet = mergeVariables(passedInVariables, session.Variables, mergedProfiles.Variables)

	// Replace variables again in the URL, in case some extra variables from session were loaded
	configuredRequest.URL = variables.ReplaceVariables(configuredRequest.URL, finalVariableSet)
	configuredRequest.Body = variables.ReplaceVariables(configuredRequest.Body, finalVariableSet)
	configuredRequest.Headers = replaceVariablesInMapOfArrayOfStrings(configuredRequest.Headers, finalVariableSet)
	configuredRequest.QueryParams = replaceVariablesInMapOfArrayOfStrings(configuredRequest.QueryParams, finalVariableSet)

	for _, cookie := range session.Cookies {
		configuredRequest.Cookies = append(configuredRequest.Cookies, cookie)
	}

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
