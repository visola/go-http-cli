package request

import (
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/variables/variables"
)

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
