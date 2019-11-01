package request

import (
	"net/url"
	"strings"

	"github.com/visola/go-http-cli/pkg/profile"
	"github.com/visola/variables/variables"
)

func replaceRequestVariables(configuredRequest Request, mergedProfiles profile.Options, context ExecutionContext) (Request, error) {
	finalVariableSet := mergeVariables(context.Variables, context.Session.Variables, mergedProfiles.Variables)

	// Replace variables again in the URL, in case some extra variables from session were loaded
	configuredRequest.URL = variables.ReplaceVariables(configuredRequest.URL, finalVariableSet)
	configuredRequest.Headers = replaceVariablesInMapOfArrayOfStrings(configuredRequest.Headers, finalVariableSet)
	configuredRequest.QueryParams = replaceVariablesInMapOfArrayOfStrings(configuredRequest.QueryParams, finalVariableSet)

	newBody, err := replaceVariablesInBody(configuredRequest, finalVariableSet)
	if err != nil {
		return configuredRequest, err
	}
	configuredRequest.Body = newBody

	for _, cookie := range context.Session.Cookies {
		configuredRequest.Cookies = append(configuredRequest.Cookies, cookie)
	}

	return configuredRequest, nil
}

func replaceVariablesInBody(configuredRequest Request, finalVariableSet map[string]string) (string, error) {
	contentType := getContentType(configuredRequest.Headers)

	// If body is form, it might have been URL encoded with variables
	// needs to be decoded, replaced, re-encoded
	if strings.HasSuffix(strings.TrimSpace(contentType), urlEncodedMimeType) {
		vals, err := url.ParseQuery(configuredRequest.Body)
		if err != nil {
			return "", err
		}

		newVals := url.Values{}
		for key, values := range vals {
			newKey := variables.ReplaceVariables(key, finalVariableSet)
			for _, value := range values {
				newVals.Add(newKey, variables.ReplaceVariables(value, finalVariableSet))
			}
		}

		return newVals.Encode(), nil
	}

	return variables.ReplaceVariables(configuredRequest.Body, finalVariableSet), nil
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
