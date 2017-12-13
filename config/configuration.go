// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

import (
	"errors"
	"os"
	"strings"

	"github.com/visola/go-http-cli/config/yaml"
	"github.com/visola/go-http-cli/options"
	"github.com/visola/go-http-cli/profile"
)

const (
	yamlExtension = ".yaml"
	ymlExtension  = ".yml"
)

// Configuration stores all the configuration that will be used to build the request.
type Configuration interface {
	BaseURL() string
	Headers() map[string][]string
	Body() string
	Method() string
	URL() string
	Variables() map[string]string
}

func hasYAMLExtension(path string) bool {
	return strings.HasSuffix(path, yamlExtension) || strings.HasSuffix(path, ymlExtension)
}

func loadProfiles(basePath string, profiles []string) ([]Configuration, error) {
	result := make([]Configuration, 0)
	for _, profile := range profiles {
		fileName := basePath + "/" + profile
		fileNameWithExtension := fileName
		if !hasYAMLExtension(fileName) {
			fileNameWithExtension = fileName + ymlExtension
			// Check for file with .yml extension
			if _, err := os.Stat(fileNameWithExtension); os.IsNotExist(err) {
				// Otherwise use .yaml
				fileNameWithExtension = fileName + yamlExtension
				// If file still doesn't exist
				if _, err := os.Stat(fileNameWithExtension); os.IsNotExist(err) {
					return result, errors.New("Configuration file does not exist: " + fileNameWithExtension)
				}
			}
		}

		yamlConfig, err := yaml.ReadFrom(fileNameWithExtension)
		if err != nil {
			return result, err
		}
		result = append(result, yamlConfig)
	}
	return result, nil
}

// Parse parses arguments and create a Configuration object.
func Parse(options *options.CommandLineOptions) (Configuration, error) {
	configurations := []Configuration{}

	if len(options.Profiles) > 0 {
		profilesDir, err := profile.GetProfilesDir()
		if err != nil {
			return nil, err
		}
		yamlConfigurations, err := loadProfiles(profilesDir, options.Profiles)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, yamlConfigurations...)
	}

	configurations = append(configurations, &BasicConfiguration{
		BodyField:      options.Body,
		HeadersField:   options.Headers,
		MethodField:    options.Method,
		URLField:       options.URL,
		VariablesField: options.Variables,
	})

	result := hierarchicalConfigurationFormat{
		configurations: configurations,
	}
	return result, nil
}
