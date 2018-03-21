package profile

import (
	"github.com/visola/go-http-cli/authorization"
)

// Used to unmarshal data from YAML files
type yamlProfileFormat struct {
	Auth      authConfiguration `yaml:"auth"`
	BaseURL   string            `yaml:"baseURL"`
	Headers   map[string]arrayOrString
	Import    arrayOrString `yaml:"import"`
	Variables map[string]string
	Requests  map[string]requestConfiguration
}

// Used to unmarshal auth options from yaml files
type authConfiguration struct {
	AuthType string `yaml:"type"`
	Password string
	Token    string
	Username string
}

// Used to unmarshal request options from yaml files
type requestConfiguration struct {
	Body         string
	FileToUpload string `yaml:"fileToUpload"`
	Headers      map[string]arrayOrString
	Method       string
	URL          string
	Values       map[string]arrayOrString
}

func (loadedProfile yamlProfileFormat) toOptions() (*Options, error) {
	headers, headersError := generateHeaders(loadedProfile)

	if headersError != nil {
		return nil, headersError
	}

	return &Options{
		BaseURL:        loadedProfile.BaseURL,
		Headers:        headers,
		RequestOptions: toMapOfRequestOptions(loadedProfile.Requests),
		Variables:      loadedProfile.Variables,
	}, nil
}

func generateHeaders(loadedProfile yamlProfileFormat) (map[string][]string, error) {
	result := toMapOfArrayOfStrings(loadedProfile.Headers)

	if loadedProfile.Auth.AuthType != "" {
		loadedAuth := loadedProfile.Auth
		auth := authorization.Authorization{
			AuthorizationType: loadedAuth.AuthType,
			Password:          loadedAuth.Password,
			Username:          loadedAuth.Username,
			Token:             loadedAuth.Token,
		}

		authValue, authError := auth.ToHeaderValue()
		if authError != nil {
			return nil, authError
		}

		result[auth.ToHeaderKey()] = []string{authValue}
	}

	return result, nil
}

func toMapOfRequestOptions(requestConfigurations map[string]requestConfiguration) map[string]RequestOptions {
	result := make(map[string]RequestOptions)
	for name, requestConfiguration := range requestConfigurations {
		result[name] = RequestOptions{
			Body:         requestConfiguration.Body,
			FileToUpload: requestConfiguration.FileToUpload,
			Headers:      toMapOfArrayOfStrings(requestConfiguration.Headers),
			Method:       requestConfiguration.Method,
			URL:          requestConfiguration.URL,
			Values:       toMapOfArrayOfStrings(requestConfiguration.Values),
		}
	}

	return result
}
