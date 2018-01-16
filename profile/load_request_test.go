package profile

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRequestOptions(t *testing.T) {
	t.Run("Loads the correct request", testLoadsCorrectRequest)
	t.Run("Returns error if request is not found", testErrorOnRequestNotFound)
}

func testLoadsCorrectRequest(t *testing.T) {
	tempProfilesDir := setupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "requests:\n  myRequest:\n    url: some/path\n  anotherRequest:\n    url: another/path\n"
	createProfile(profileName, profileContent, tempProfilesDir)

	request, requestErr := LoadRequestOptions("myRequest", []string{profileName})

	assert.Nil(t, requestErr, "Should load request correctly")
	if requestErr != nil {
		panic(requestErr)
	}

	assert.Equal(t, "some/path", request.URL, "Should load the correct URL")
}

func testErrorOnRequestNotFound(t *testing.T) {
	tempProfilesDir := setupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "requests:\n  myRequest:\n    url: some/path\n"
	createProfile(profileName, profileContent, tempProfilesDir)

	_, requestErr := LoadRequestOptions("anotherRequest", []string{profileName})

	assert.NotNil(t, requestErr, "Should fail to load if request name not found")
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
