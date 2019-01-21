package profile

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// GetAvailableProfiles return the name of all profiles available in the configured directory
func GetAvailableProfiles() ([]string, error) {
	profilesDirectory, getDirErr := GetProfilesDir()

	if getDirErr != nil {
		return nil, getDirErr
	}

	files, readDirErr := ioutil.ReadDir(profilesDirectory)

	if readDirErr != nil {
		return nil, readDirErr
	}

	result := make([]string, 0)

	for _, file := range files {
		extension := filepath.Ext(file.Name())
		if extension == ".yml" || extension == ".yaml" {
			result = append(result, strings.TrimSuffix(file.Name(), extension))
		}
	}

	return result, nil
}
