// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

import (
	"errors"
	"os"
	"strings"

	"github.com/visola/go-http-cli/config/yaml"
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

func loadConfigurations(paths []string) ([]Configuration, error) {
	result := make([]Configuration, 0)
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return result, errors.New("Configuration file does not exist: " + path)
		}
		yamlConfig, err := yaml.ReadFrom(path)
		if err != nil {
			return result, err
		}
		result = append(result, yamlConfig)
	}
	return result, nil
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
func Parse(args []string) (Configuration, error) {
	commandLineConfiguration, err := parseCommandLine(args)
	if err != nil {
		return nil, err
	}

	configurations := []Configuration{}

	if len(commandLineConfiguration.profiles) > 0 {
		profilesDir, err := profile.GetProfilesDir()
		if err != nil {
			return nil, err
		}
		yamlConfigurations, err := loadProfiles(profilesDir, commandLineConfiguration.profiles)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, yamlConfigurations...)
	}

	if len(commandLineConfiguration.configurationPaths) > 0 {
		yamlConfigurations, err := loadConfigurations(commandLineConfiguration.configurationPaths)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, yamlConfigurations...)
	}

	configurations = append(configurations, commandLineConfiguration)

	result := hierarchicalConfigurationFormat{
		configurations: configurations,
	}
	return result, nil
}
