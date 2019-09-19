package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/visola/go-http-cli/pkg/base"
	"github.com/visola/go-http-cli/pkg/profile"
)

const jsonMimeType = "application/json"
const urlEncodedMimeType = "application/x-www-form-urlencoded"

var bodyBuilderContentTypes = [...]string{
	urlEncodedMimeType,
	jsonMimeType,
}

// CreateConfigureRequestOptions creates the options to configure a request based on the options
func CreateConfigureRequestOptions(passedInOptions ...ConfigureRequestOption) *ConfigureRequestOptions {
	configureOptions := &ConfigureRequestOptions{}
	for _, configureOption := range passedInOptions {
		configureOption(configureOptions)
	}

	return configureOptions
}

// ConfigureRequestSimple is a simpler version of configure request
// TODO - This was brute forced here to solve the post-processing issue, this code needs to be refactored
func ConfigureRequestSimple(unconfiguredRequest Request, mergedProfile *profile.Options, requestName string) (*Request, error) {
	namedRequest, namedRequestErr := profile.FindNamedRequest(mergedProfile, requestName)
	if namedRequestErr != nil {
		return nil, namedRequestErr
	}

	configuredRequest := Request{}
	configuredRequest.Merge(mergedProfile)
	configuredRequest.Merge(namedRequest)
	configuredRequest.Merge(unconfiguredRequest)

	finalValueSet := getValues(namedRequest)

	return finalizeConfiguringRequest(configuredRequest, mergedProfile, namedRequest, finalValueSet)
}

// ConfigureRequest configures a request to be executed based on the provided options
func ConfigureRequest(unconfiguredRequest Request, mergedProfile *profile.Options, configureOptions *ConfigureRequestOptions) (*Request, error) {
	namedRequest, namedRequestErr := profile.FindNamedRequest(mergedProfile, configureOptions.RequestName)
	if namedRequestErr != nil {
		return nil, namedRequestErr
	}

	configuredRequest := Request{}
	configuredRequest.Merge(mergedProfile)
	configuredRequest.Merge(namedRequest)
	configuredRequest.Merge(unconfiguredRequest)
	finalValueSet := getValues(namedRequest, configureOptions)

	return finalizeConfiguringRequest(configuredRequest, mergedProfile, namedRequest, finalValueSet)
}

func finalizeConfiguringRequest(configuredRequest Request, mergedProfile *profile.Options, namedRequest profile.NamedRequest, finalValueSet map[string][]string) (*Request, error) {
	hasBody := configuredRequest.Body != ""
	hasValues := len(finalValueSet) > 0
	hasContentType := getContentType(configuredRequest.Headers) != ""

	configuredRequest.Method = getMethod(configuredRequest.Method, hasBody)

	if !hasContentType && (hasBody || (hasValues && configuredRequest.Method != http.MethodGet)) {
		configuredRequest.Headers["Content-Type"] = []string{jsonMimeType}
	}

	var createdFromValues bool
	configuredRequest.Body, createdFromValues = getBody(configuredRequest, finalValueSet)
	configuredRequest.URL = ParseURL(mergedProfile.BaseURL, configuredRequest.URL, namedRequest.URL)

	if !createdFromValues && len(finalValueSet) > 0 {
		configuredRequest.QueryParams = finalValueSet
	}

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
		// TODO - Fix this for variables passed in values or keys
		return encodeValues(values)
	}

	return fmt.Sprintf("Unsupported body type: %s", contentType)
}

func getBody(configuredRequest Request, values map[string][]string) (string, bool) {
	if configuredRequest.Method == http.MethodGet {
		return "", false
	}

	if configuredRequest.Body != "" {
		return configuredRequest.Body, false
	}

	if len(values) > 0 {
		return createBody(configuredRequest, values), true
	}

	return "", false
}

func getMethod(currentMethod string, hasBody bool) string {
	method := currentMethod

	if method != "" {
		return method
	}

	if hasBody {
		return http.MethodPost
	}

	return http.MethodGet
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
