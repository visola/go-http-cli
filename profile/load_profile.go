package profile

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	yamlExtension = ".yaml"
	ymlExtension  = ".yml"
)

// LoadProfile loads Options for a specific profile by name.
func LoadProfile(profileName string) (*Options, error) {
	profilesDir, profilesDirErr := GetProfilesDir()
	if profilesDirErr != nil {
		return nil, profilesDirErr
	}

	fileName := profilesDir + "/" + profileName
	fileNameWithExtension := fileName
	if !hasYAMLExtension(fileName) {
		fileNameWithExtension = fileName + ymlExtension
		// Check for file with .yml extension
		if _, err := os.Stat(fileNameWithExtension); os.IsNotExist(err) {
			// Otherwise use .yaml
			fileNameWithExtension = fileName + yamlExtension
			// If file still doesn't exist
			if _, err := os.Stat(fileNameWithExtension); os.IsNotExist(err) {
				return nil, errors.New("Configuration file does not exist: " + fileNameWithExtension)
			}
		}
	}

	return readFrom(fileNameWithExtension)
}

func readFrom(pathToYamlFile string) (*Options, error) {
	loadedOptions := new(yamlProfileFormat)

	var err error
	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(pathToYamlFile)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedOptions)

	if err != nil {
		return nil, err
	}

	return loadedOptions.toOptions(), nil
}

func hasYAMLExtension(path string) bool {
	return strings.HasSuffix(path, yamlExtension) || strings.HasSuffix(path, ymlExtension)
}
