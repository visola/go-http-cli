package config

import (
	"errors"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type keyValuePair []string

func (i *keyValuePair) String() string {
	return ""
}

func (i *keyValuePair) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *keyValuePair) Type() string {
	return "keyValuePair"
}

type commandLineConfiguration struct {
	body               string
	configurationPaths []string
	headers            map[string][]string
	method             string
	profiles           []string
	url                string
	variables          map[string]string
}

func (conf commandLineConfiguration) BaseURL() string {
	return ""
}

func (conf commandLineConfiguration) Headers() map[string][]string {
	return conf.headers
}

func (conf commandLineConfiguration) Body() string {
	return conf.body
}

func (conf commandLineConfiguration) Method() string {
	return conf.method
}

func (conf commandLineConfiguration) URL() string {
	return conf.url
}

func (conf commandLineConfiguration) Variables() map[string]string {
	return conf.variables
}

func parseMultiValues(headers keyValuePair) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, kv := range headers {
		s := strings.Split(kv, "=")
		if len(s) != 2 {
			return result, errors.New("Error while parsing key value pair '" + kv + "'\nShould be an '=' separated key/value, e.g.: Content-type=application/x-www-form-urlencoded")
		}

		key := s[0]
		value := s[1]

		if existingValue, ok := result[key]; ok {
			result[key] = append(existingValue, value)
		} else {
			result[key] = []string{value}
		}
	}

	return result, nil
}

func parseValues(headers keyValuePair) (map[string]string, error) {
	result := make(map[string]string)

	for _, kv := range headers {
		s := strings.Split(kv, "=")
		if len(s) != 2 {
			return result, errors.New("Error while parsing key value pair '" + kv + "'\nShould be an '=' separated key/value, e.g.: Content-type=application/x-www-form-urlencoded")
		}

		key := s[0]
		value := s[1]

		result[key] = value
	}

	return result, nil
}

func parseArgs(args []string) (string, []string, error) {
	if len(args) == 0 {
		return "", nil, errors.New("no arguments passed in")
	}

	profiles := make([]string, 0)
	var url string
	for _, arg := range args {
		if arg[0] == '+' { // Found a profile to activate
			profiles = append(profiles, arg[1:])
		} else {
			url = arg
		}
	}
	return url, profiles, nil
}

func parseCommandLine(args []string) (*commandLineConfiguration, error) {
	var method string
	var body string
	var headers keyValuePair
	var configPaths keyValuePair
	var variables keyValuePair

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	commandLine.StringVarP(&method, "method", "X", "", "HTTP method to be used")
	commandLine.StringVarP(&body, "data", "d", "", "Data to be sent as body")
	commandLine.VarP(&headers, "header", "H", "Headers to include with your request")
	commandLine.VarP(&configPaths, "config", "c", "Path to configuration files to be used")
	commandLine.VarP(&variables, "variable", "V", "Variables to be used on substitutions")

	commandLine.Parse(args)

	if method == "" {
		if body == "" {
			method = "GET"
		} else {
			method = "POST"
		}
	}

	result := new(commandLineConfiguration)
	result.method = method
	result.body = body

	result.configurationPaths = configPaths

	url, profiles, urlError := parseArgs(commandLine.Args())
	result.url = url
	result.profiles = profiles

	if urlError != nil {
		return result, urlError
	}

	parsedHeaders, headerError := parseMultiValues(headers)
	result.headers = parsedHeaders

	if headerError != nil {
		return result, headerError
	}

	parsedVariables, variableError := parseValues(variables)
	result.variables = parsedVariables

	if variableError != nil {
		return result, variableError
	}

	return result, nil
}
