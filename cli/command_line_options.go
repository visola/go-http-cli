package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// CommandLineOptions stores information that was requested by the user from the CLI.
type CommandLineOptions struct {
	Body            string
	Headers         map[string][]string
	FollowLocation  bool
	FileToUpload    string
	MaxRedirect     int
	Method          string
	OutputFile      string
	PostProcessFile string
	Profiles        []string
	RequestName     string
	URL             string
	Values          map[string][]string
	Variables       map[string]string
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
	var body, fileToUpload, method, outputFile, postProcessFile string
	var configPaths, headers, variables keyValuePair
	var followLocation bool

	commandLine := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	commandLine.VarP(&configPaths, "config", "c", "Path to configuration files to be used")
	commandLine.StringVarP(&body, "data", "d", "", "Data to be sent as body")
	commandLine.VarP(&headers, "header", "H", "Headers to include with your request")
	commandLine.BoolVarP(&followLocation, "location", "L", false, "Automatically follow redirects")
	maxRedirect := commandLine.Int("max-redirs", 50, "Maximum number of redirects to follow")
	commandLine.StringVarP(&method, "method", "X", "", "HTTP method to be used")
	commandLine.StringVarP(&outputFile, "output", "o", "", "File to save the response")
	commandLine.StringVarP(&postProcessFile, "post-process", "", "", "Javascript file to post process the request/response")
	commandLine.StringVarP(&fileToUpload, "upload-file", "T", "", "Path to the file to be uploaded")
	commandLine.VarP(&variables, "variable", "V", "Variables to be used on substitutions")

	commandLine.Parse(args)

	result := new(CommandLineOptions)

	result.Body = body
	result.FileToUpload = fileToUpload
	result.FollowLocation = followLocation
	result.MaxRedirect = *maxRedirect
	result.Method = method
	result.OutputFile = outputFile
	result.PostProcessFile = postProcessFile

	url, requestName, profiles, values, urlError := parseArgs(commandLine.Args())
	result.Profiles = profiles
	result.RequestName = requestName
	result.URL = url
	result.Values = values

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

func extractKeyValuePair(arg string) (string, string) {
	indexToSplit := getIndexToSplit(arg)
	questionMarkIndex := strings.Index(arg, "?")

	// If no separator found or if there's a question mark before the separator
	// question mark before the separator, for now assume it's a partial URL
	if indexToSplit < 0 || (questionMarkIndex >= 0 && questionMarkIndex < indexToSplit) {
		return "", ""
	}

	key := arg[0:indexToSplit]
	value := arg[indexToSplit+1:]
	if value == "" {
		value = fmt.Sprintf("{%s}", key)
	}

	return key, value
}

func getIndexToSplit(keyValuePair string) int {
	equalIndex := strings.Index(keyValuePair, "=")
	colonIndex := strings.Index(keyValuePair, ":")
	indexToSplit := -1

	if equalIndex > 0 {
		indexToSplit = equalIndex
	}

	if colonIndex > 0 && (indexToSplit == -1 || colonIndex < indexToSplit) {
		indexToSplit = colonIndex
	}

	return indexToSplit
}

func parseArgs(args []string) (string, string, []string, map[string][]string, error) {
	if len(args) == 0 {
		return "", "", nil, nil, errors.New("no arguments passed in")
	}

	profiles := make([]string, 0)
	values := make(map[string][]string)
	var requestName, url string
	for _, arg := range args {
		key, value := extractKeyValuePair(arg)
		if arg[0] == '+' { // Found a profile to activate
			profiles = append(profiles, arg[1:])
		} else if arg[0] == '@' {
			requestName = arg[1:]
		} else if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			url = arg
		} else if key != "" {
			if existingValue, ok := values[key]; ok {
				values[key] = append(existingValue, value)
			} else {
				values[key] = []string{value}
			}
		} else {
			url = arg
		}
	}
	return url, requestName, profiles, values, nil
}

func parseMultiValues(headers keyValuePair) (map[string][]string, error) {
	result := make(map[string][]string)

	for _, kv := range headers {
		indexToSplit := getIndexToSplit(kv)
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
