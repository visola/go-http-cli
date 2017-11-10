package profile

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetProfilesDir return the directory where profiles are stored
func GetProfilesDir() (string, error) {
	profilesDir := os.Getenv("GO_HTTP_PROFILES")
	if profilesDir == "" {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		profilesDir = user.HomeDir + "/go-http-cli"
	}
	return profilesDir, nil
}

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
