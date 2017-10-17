package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// We want to accept headers as single string or array of strings
type headerValue []string

func (v *headerValue) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		// Apparently value is not an array
		var single string
		err := unmarshal(&single)
		if err != nil {
			// Still can't parse it
			return err
		}
		*v = []string{single}
	} else {
		*v = multi
	}
	return nil
}

// Used to unmarshal data from YAML files
type yamlConfigurationFormat struct {
	Headers map[string]headerValue
}

// Configuration implementation that wraps configuration coming from a YAML file
type fileConfiguration struct {
	parsedYaml *yamlConfigurationFormat
}

func (conf fileConfiguration) Headers() map[string][]string {
	result := make(map[string][]string)
	for header, values := range conf.parsedYaml.Headers {
		result[header] = values
	}
	return result
}

func (conf fileConfiguration) Body() string {
	return ""
}

func (conf fileConfiguration) Method() string {
	return ""
}

func (conf fileConfiguration) URL() string {
	return ""
}

func readFrom(pathToYamlFile string) (*fileConfiguration, error) {
	yamlConfiguration := new(yamlConfigurationFormat)

	var err error
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(pathToYamlFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &yamlConfiguration)

	if err != nil {
		return nil, err
	}

	result := &fileConfiguration{parsedYaml: yamlConfiguration}

	return result, nil
}
