// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

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
	result := hierarchicalConfigurationFormat{
		configurations: []Configuration{commandLineConfiguration},
	}
	return result, nil
}
