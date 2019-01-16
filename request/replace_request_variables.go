package request

import (
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/variables/variables"
)

func replaceRequestVariables(configuredRequest Request, mergedProfiles profile.Options, context ExecutionContext) (Request, error) {
	finalVariableSet := mergeVariables(context.Variables, context.Session.Variables, mergedProfiles.Variables)

	// Replace variables again in the URL, in case some extra variables from session were loaded
	configuredRequest.URL = variables.ReplaceVariables(configuredRequest.URL, finalVariableSet)
	configuredRequest.Body = variables.ReplaceVariables(configuredRequest.Body, finalVariableSet)
	configuredRequest.Headers = replaceVariablesInMapOfArrayOfStrings(configuredRequest.Headers, finalVariableSet)
	configuredRequest.QueryParams = replaceVariablesInMapOfArrayOfStrings(configuredRequest.QueryParams, finalVariableSet)

	for _, cookie := range context.Session.Cookies {
		configuredRequest.Cookies = append(configuredRequest.Cookies, cookie)
	}

	return configuredRequest, nil
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
