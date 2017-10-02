package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

type headerFlags []string

func (i *headerFlags) String() string {
	return "No String Representation"
}

func (i *headerFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Configuration struct {
	Body    string
	Headers map[string]string
	Method  string
	Url     string
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

func parseUrl(args []string) (string, error) {
	return args[0], nil
}

func Parse(args []string) (*Configuration, error) {
	var method string
	var body string
	var headers headerFlags

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	commandLine.StringVar(&method, "method", "GET", "HTTP method to be used")
	commandLine.StringVar(&body, "data", "", "Data to be sent as body")
	commandLine.Var(&headers, "header", "Headers to include with your request")

	commandLine.Parse(args)

	result := new(Configuration)
	result.Method = method
	result.Body = body

	url, urlError := parseUrl(commandLine.Args())
	result.Url = url

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
