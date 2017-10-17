// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

import (
	"errors"
	"os"
	"os/user"
	"strings"
)

const (
	yamlExtension = ".yaml"
	ymlExtension  = ".yml"
)

// Configuration stores all the configuration that will be used to build the request.
type Configuration interface {
	Headers() map[string][]string
	Body() string
	Method() string
	URL() string
}

func hasYAMLExtension(path string) bool {
	return strings.HasSuffix(path, yamlExtension) || strings.HasSuffix(path, ymlExtension)
}

// Get the directory where the user store his profiles
func getProfilesDir() (string, error) {
	profilesDir := os.Getenv("GO_HTTP_PROFILES")
	if profilesDir == "" {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		profilesDir = user.HomeDir + "/go-http-cli"
	}
	return profilesDir, nil
}

func loadConfigurations(paths []string) ([]Configuration, error) {
	result := make([]Configuration, 0)
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return result, errors.New("Configuration file does not exist: " + path)
		}
		yamlConfig, err := readFrom(path)
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

		yamlConfig, err := readFrom(fileNameWithExtension)
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
		profilesDir, err := getProfilesDir()
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
