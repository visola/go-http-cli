// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

import (
	"errors"
	"os"
)

// Configuration stores all the configuration that will be used to build the request.
type Configuration interface {
	Headers() map[string][]string
	Body() string
	Method() string
	URL() string
}

// Parse parses arguments and create a Configuration object.
func Parse(args []string) (Configuration, error) {
	commandLineConfiguration, err := parseCommandLine(args)
	if err != nil {
		return nil, err
	}

	// We'll have at least one configuration
	configurations := []Configuration{}

	if len(commandLineConfiguration.ConfigurationPaths) > 0 {
		for _, configPath := range commandLineConfiguration.ConfigurationPaths {
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return commandLineConfiguration, errors.New("Configuration file does not exist: " + configPath)
			}
			yamlConfig, err := readFrom(configPath)
			if err != nil {
				return commandLineConfiguration, err
			}
			configurations = append(configurations, yamlConfig)
		}
	}

	configurations = append(configurations, commandLineConfiguration)

	result := hierarchicalConfigurationFormat{
		configurations: configurations,
	}
	return result, nil
}
