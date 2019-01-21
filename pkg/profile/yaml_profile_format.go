package profile

import (
	"github.com/visola/go-http-cli/pkg/authorization"
	"github.com/visola/go-http-cli/pkg/model"
)

// Used to unmarshal data from YAML files
type yamlProfileFormat struct {
	Auth      authConfiguration `yaml:"auth"`
	BaseURL   string            `yaml:"baseURL"`
	Headers   map[string]model.ArrayOrString
	Import    model.ArrayOrString `yaml:"import"`
	Requests  map[string]requestConfiguration
	Variables map[string]string
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
	Body              string
	FileToUpload      string `yaml:"fileToUpload"`
	Headers           map[string]model.ArrayOrString
	Method            string
	PostProcessScript string `yaml:"postProcessScript"`
	URL               string
	Values            map[string]model.ArrayOrString
}

func (loadedProfile yamlProfileFormat) toOptions() (*Options, error) {
	headers, headersError := generateHeaders(loadedProfile)

	if headersError != nil {
		return nil, headersError
	}

	return &Options{
		BaseURL:      loadedProfile.BaseURL,
		Headers:      headers,
		NamedRequest: toMapOfNamedRequest(loadedProfile.Requests),
		Variables:    loadedProfile.Variables,
	}, nil
}

func generateHeaders(loadedProfile yamlProfileFormat) (map[string][]string, error) {
	result := model.ToMapOfArrayOfStrings(loadedProfile.Headers)

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

func toMapOfNamedRequest(requestConfigurations map[string]requestConfiguration) map[string]NamedRequest {
	result := make(map[string]NamedRequest)
	for name, requestConfiguration := range requestConfigurations {
		result[name] = NamedRequest{
			Body:              requestConfiguration.Body,
			FileToUpload:      requestConfiguration.FileToUpload,
			Headers:           model.ToMapOfArrayOfStrings(requestConfiguration.Headers),
			Method:            requestConfiguration.Method,
			PostProcessScript: requestConfiguration.PostProcessScript,
			URL:               requestConfiguration.URL,
			Values:            model.ToMapOfArrayOfStrings(requestConfiguration.Values),
		}
	}

	return result
}
