package yaml

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// We want to accept headers as single string or array of strings
type arrayOrString []string

func (v *arrayOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
	BaseURL   string `yaml:"baseURL"`
	Headers   map[string]arrayOrString
	Variables map[string]arrayOrString
}

// FileConfiguration represents a Configuration loaded from a YAML file
type FileConfiguration struct {
	parsedYaml *yamlConfigurationFormat
}

// BaseURL returns the base URL loaded from the file
func (conf FileConfiguration) BaseURL() string {
	return conf.parsedYaml.BaseURL
}

// Headers loaded from the file
func (conf FileConfiguration) Headers() map[string][]string {
	result := make(map[string][]string)
	for header, values := range conf.parsedYaml.Headers {
		result[header] = values
	}
	return result
}

// Body returns empty string (not implemented)
func (conf FileConfiguration) Body() string {
	return ""
}

// Method returns empty string (not implemented)
func (conf FileConfiguration) Method() string {
	return ""
}

// URL returns empty string  (not implemented)
func (conf FileConfiguration) URL() string {
	return ""
}

// Variables that were added to the configuration file
func (conf FileConfiguration) Variables() map[string][]string {
	result := make(map[string][]string)
	for name, values := range conf.parsedYaml.Variables {
		result[name] = values
	}
	return result
}

// ReadFrom a YAML file and creates a configuration
func ReadFrom(pathToYamlFile string) (*FileConfiguration, error) {
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

	result := &FileConfiguration{parsedYaml: yamlConfiguration}

	return result, nil
}
