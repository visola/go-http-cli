// Package config has all the things required to parse configuration from command line arguments
// and files.
package config

import (
	"errors"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type headerFlags []string

func (i *headerFlags) String() string {
	return "No String Representation"
}

func (i *headerFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *headerFlags) Type() string {
	return "headers"
}

// Configuration stores all the configuration that will be used to build the request.
type Configuration struct {
	Body    string
	Headers map[string]string
	Method  string
	URL     string
}

func parseHeaders(headers headerFlags) (map[string]string, error) {
	result := make(map[string]string)

	for _, kv := range headers {
		s := strings.Split(kv, "=")
		if len(s) != 2 {
			return result, errors.New("Error while parsing header '" + kv + "'\nShould be a '=' separated key/value, e.g.: Content-type=application/x-www-form-urlencoded")
		}
		result[s[0]] = s[1]
	}

	return result, nil
}

func parseURL(args []string) (string, error) {
	return args[0], nil
}

// Parse parses arguments and create a Configuration object.
func Parse(args []string) (*Configuration, error) {
	var method string
	var body string
	var headers headerFlags

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	commandLine.StringVarP(&method, "method", "X", "GET", "HTTP method to be used")
	commandLine.StringVarP(&body, "data", "d", "", "Data to be sent as body")
	commandLine.VarP(&headers, "header", "H", "Headers to include with your request")

	commandLine.Parse(args)

	result := new(Configuration)
	result.Method = method
	result.Body = body

	url, urlError := parseURL(commandLine.Args())
	result.URL = url

	if urlError != nil {
		return result, urlError
	}

	parsedHeaders, headerError := parseHeaders(headers)
	result.Headers = parsedHeaders

	if headerError != nil {
		return result, headerError
	}

	return result, nil
}
