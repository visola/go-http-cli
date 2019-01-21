package profile

import (
	"fmt"
	"io/ioutil"
	"os"
)

// CreateTestProfile helper method for testing. Will create a profile file with the specified
// content in the specified directory.
func CreateTestProfile(profileName string, profileContent string, profilesDir string) {
	profileFile, err := os.Create(fmt.Sprintf("%s/%s.yml", profilesDir, profileName))
	if err != nil {
		panic(err)
	}

	defer profileFile.Close()

	profileFile.WriteString(profileContent)
}

// SetupTestProfilesDir helper method for testing. This will create a temporary directory where
// profiles can be dropped in. It will also set the environment variable so that profile files
// are read from the temporary directory.
func SetupTestProfilesDir() string {
	dir, err := ioutil.TempDir("", "profiles")
	if err != nil {
		panic(err)
	}

	os.Setenv(profilesDirEnvVariable, dir)
	return dir
}
