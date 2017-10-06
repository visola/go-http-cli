package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Used to unmarshal data from YAML files
type yamlConfigurationFormat struct {
	Headers map[string][]string
}

// Configuration implementation that wraps configuration coming from a YAML file
type fileConfiguration struct {
	parsedYaml *yamlConfigurationFormat
}

func (conf fileConfiguration) Headers() map[string][]string {
	return conf.parsedYaml.Headers
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
