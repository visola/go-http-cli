package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/visola/go-http-cli/base"
	"github.com/visola/go-http-cli/profile"
)

const jsonMimeType = "application/json"
const urlEncodedMimeType = "application/x-www-form-urlencoded"

var bodyBuilderContentTypes = [...]string{
	urlEncodedMimeType,
	jsonMimeType,
}

// ConfigureRequest configures a request to be executed based on the provided options
func ConfigureRequest(unconfiguredRequest Request, passedInOptions ...ConfigureRequestOption) (*Request, error) {
	configureOptions := &ConfigureRequestOptions{}
	for _, configureOption := range passedInOptions {
		configureOption(configureOptions)
	}

	mergedProfile, profileError := profile.LoadAndMergeProfiles(configureOptions.ProfileNames)
	if profileError != nil {
		return nil, profileError
	}

	namedRequest, namedRequestErr := findNamedRequest(mergedProfile, configureOptions.RequestName)
	if namedRequestErr != nil {
		return nil, namedRequestErr
	}

	configuredRequest := unconfiguredRequest

	configuredRequest.Merge(mergedProfile)
	configuredRequest.Merge(namedRequest)
	configuredRequest.Merge(unconfiguredRequest)
	finalValueSet := getValues(namedRequest, configureOptions)

	hasBody := configuredRequest.Body != ""
	hasValues := len(finalValueSet) > 0
	hasContentType := getContentType(configuredRequest.Headers) != ""

	if !hasContentType && (hasBody || hasValues) {
		configuredRequest.Headers["Content-Type"] = []string{jsonMimeType}
	}

	configuredRequest.Method = getMethod(configuredRequest.Method, getContentType(configuredRequest.Headers), hasBody, hasValues)

	configuredRequest.Body = getBody(configuredRequest, finalValueSet)
	configuredRequest.URL = ParseURL(mergedProfile.BaseURL, configuredRequest.URL, namedRequest.URL)

	return &configuredRequest, nil
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

func createBody(processedRequest Request, values map[string][]string) string {
	contentType := getContentType(processedRequest.Headers)
	if contentType == "" || strings.HasSuffix(strings.TrimSpace(contentType), jsonMimeType) {
		return buildJSON(values)
	} else if strings.HasSuffix(strings.TrimSpace(contentType), urlEncodedMimeType) {
		return encodeValues(values)
	}

	return fmt.Sprintf("Unsupported body type: %s", contentType)
}

func findNamedRequest(mergedProfile profile.Options, requestName string) (profile.NamedRequest, error) {
	if requestName == "" {
		return profile.NamedRequest{}, nil
	}

	var namedRequest profile.NamedRequest

	var exists bool
	if namedRequest, exists = mergedProfile.NamedRequest[requestName]; requestName != "" && !exists {
		return profile.NamedRequest{}, fmt.Errorf("Request with name %s not found", requestName)
	}

	return namedRequest, nil
}

func getBody(configuredRequest Request, values map[string][]string) string {
	if configuredRequest.Method == http.MethodGet {
		return ""
	}

	if configuredRequest.Body != "" {
		return configuredRequest.Body
	}

	if len(values) > 0 {
		return createBody(configuredRequest, values)
	}

	return ""
}

func getMethod(currentMethod string, contentType string, hasBody bool, hasValues bool) string {
	method := currentMethod

	if method == "" {
		if hasBody {
			method = http.MethodPost
		} else {
			method = http.MethodGet
		}
	} else {
		// User knows what he wants
		return method
	}

	// If still empty
	if method == "" || method == http.MethodGet {
		// If there are values, check if they should go in the body
		if hasValues {
			for _, knownType := range bodyBuilderContentTypes {
				if strings.HasPrefix(contentType, knownType) {
					return http.MethodPost
				}
			}
		}
	}

	return method
}

func getValues(valueSets ...base.WithValues) map[string][]string {
	result := make(map[string][]string)
	for _, valueSet := range valueSets {
		for valName, val := range valueSet.GetValues() {
			result[valName] = val
		}
	}
	return result
}
