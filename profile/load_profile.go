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

// LoadAndMergeProfiles loads all profiles and merge them in the order passed in
func LoadAndMergeProfiles(profileNames []string) (Options, error) {
	profiles := make([]Options, len(profileNames))

	for index, profileName := range profileNames {
		loadedProfile, profileError := LoadProfile(profileName)
		if profileError != nil {
			return Options{}, profileError
		}
		profiles[index] = loadedProfile
	}

	return MergeOptions(profiles), nil
}

// LoadProfile loads Options for a specific profile by name.
func LoadProfile(profileName string) (loadedOptions Options, err error) {
	profilesDir, profilesDirErr := GetProfilesDir()
	if profilesDirErr != nil {
		return loadedOptions, profilesDirErr
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
				return loadedOptions, errors.New("Configuration file does not exist: " + fileNameWithExtension)
			}
		}
	}

	return readFrom(fileNameWithExtension)
}

func readFrom(pathToYamlFile string) (finalOptions Options, err error) {
	loadedOptions := new(yamlProfileFormat)

	var yamlContent []byte
	yamlContent, err = ioutil.ReadFile(pathToYamlFile)

	if err != nil {
		return finalOptions, err
	}

	err = yaml.Unmarshal(yamlContent, &loadedOptions)

	if err != nil {
		return finalOptions, err
	}

	importedOptions := make([]Options, 0)
	for _, toImport := range loadedOptions.Import {
		imported, importErr := LoadProfile(toImport)
		if importErr != nil {
			return finalOptions, err
		}
		importedOptions = append(importedOptions, imported)
	}

	readOption, conversionErr := loadedOptions.toOptions()
	if conversionErr != nil {
		return finalOptions, conversionErr
	}

	importedOptions = append(importedOptions, *readOption)
	return MergeOptions(importedOptions), err
}

func hasYAMLExtension(path string) bool {
	return strings.HasSuffix(path, yamlExtension) || strings.HasSuffix(path, ymlExtension)
}
