package cli

import (
	"errors"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// CommandLineOptions stores information that was requested by the user from the CLI.
type CommandLineOptions struct {
	Body           string
	Headers        map[string][]string
	FollowLocation bool
	FileToUpload   string
	MaxRedirect    int
	Method         string
	Profiles       []string
	RequestName    string
	URL            string
	Variables      map[string]string
}

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

// ParseCommandLineOptions parses the arguments received on the command line and generate a basic configuration.
func ParseCommandLineOptions(args []string) (*CommandLineOptions, error) {
	var method string
	var body string
	var fileToUpload string
	var headers keyValuePair
	var configPaths keyValuePair
	var variables keyValuePair
	var followLocation bool

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	maxRedirect := commandLine.Int("max-redirs", 50, "Maximum number of redirects to follow")
	commandLine.BoolVarP(&followLocation, "location", "L", false, "Automatically follow redirects")
	commandLine.StringVarP(&method, "method", "X", "", "HTTP method to be used")
	commandLine.StringVarP(&body, "data", "d", "", "Data to be sent as body")
	commandLine.StringVarP(&fileToUpload, "upload-file", "T", "", "Path to the file to be uploaded")
	commandLine.VarP(&headers, "header", "H", "Headers to include with your request")
	commandLine.VarP(&configPaths, "config", "c", "Path to configuration files to be used")
	commandLine.VarP(&variables, "variable", "V", "Variables to be used on substitutions")

	commandLine.Parse(args)

	result := new(CommandLineOptions)

	result.FollowLocation = followLocation
	result.MaxRedirect = *maxRedirect
	result.Method = method
	result.Body = body
	result.FileToUpload = fileToUpload

	url, requestName, profiles, urlError := parseArgs(commandLine.Args())
	result.URL = url
	result.RequestName = requestName
	result.Profiles = profiles

	if urlError != nil {
		return result, urlError
	}

	parsedHeaders, headerError := parseMultiValues(headers)
	result.Headers = parsedHeaders

	if headerError != nil {
		return result, headerError
	}

	parsedVariables, variableError := parseValues(variables)
	result.Variables = parsedVariables

	if variableError != nil {
		return result, variableError
	}

	return result, nil
}

func parseMultiValues(headers keyValuePair) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, kv := range headers {
		equalIndex := strings.Index(kv, "=")
		colonIndex := strings.Index(kv, ":")
		indexToSplit := -1

		if equalIndex > 0 {
			indexToSplit = equalIndex
		}

		if colonIndex > 0 && (indexToSplit == -1 || colonIndex < indexToSplit) {
			indexToSplit = colonIndex
		}

		if indexToSplit <= 0 {
			return result, errors.New("Error while parsing key value pair '" + kv + "'\nShould be an '=' or ':' separated key/value, e.g.: Content-type=application/x-www-form-urlencoded")
		}

		key := kv[0:indexToSplit]
		value := kv[indexToSplit+1:]

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

func parseArgs(args []string) (string, string, []string, error) {
	if len(args) == 0 {
		return "", "", nil, errors.New("no arguments passed in")
	}

	profiles := make([]string, 0)
	var requestName, url string
	for _, arg := range args {
		if arg[0] == '+' { // Found a profile to activate
			profiles = append(profiles, arg[1:])
		} else if arg[0] == '@' {
			requestName = arg[1:]
		} else {
			url = arg
		}
	}
	return url, requestName, profiles, nil
}
