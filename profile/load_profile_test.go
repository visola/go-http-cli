package profile

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProfile(t *testing.T) {
	t.Run("Loads basic authorization correctly", testLoadsBasicAuthorization)
}

func testLoadsBasicAuthorization(t *testing.T) {
	tempProfilesDir := setupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "auth:\n  type:  basic\n  username: myUsername\n  password: myPassword\n"
	createProfile(profileName, profileContent, tempProfilesDir)

	profile, requestErr := LoadProfile(profileName)

	assert.Nil(t, requestErr, "Should load request correctly")
	if requestErr != nil {
		panic(requestErr)
	}

	authValues, exists := profile.Headers["Authorization"]
	assert.True(t, exists, "Should have Authorization header")
	assert.Equal(t, 1, len(authValues), "Should have only one value")
}

func createProfile(profileName string, profileContent string, profilesDir string) {
	profileFile, err := os.Create(fmt.Sprintf("%s/%s.yml", profilesDir, profileName))
	if err != nil {
		panic(err)
	}

	defer profileFile.Close()

	profileFile.WriteString(profileContent)
}

func setupTestProfilesDir() string {
	dir, err := ioutil.TempDir("", "profiles")
	if err != nil {
		panic(err)
	}

	os.Setenv(profilesDirEnvVariable, dir)
	return dir
}
