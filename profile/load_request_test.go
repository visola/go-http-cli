package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRequestOptions(t *testing.T) {
	t.Run("Loads the correct request", testLoadsCorrectRequest)
	t.Run("Returns error if request is not found", testErrorOnRequestNotFound)
	t.Run("Loads basic authorization correctly", testLoadsBasicAuthorization)
}

func testLoadsCorrectRequest(t *testing.T) {
	tempProfilesDir := SetupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "requests:\n  myRequest:\n    url: some/path\n  anotherRequest:\n    url: another/path\n"
	CreateTestProfile(profileName, profileContent, tempProfilesDir)

	request, requestErr := LoadRequestOptions("myRequest", []string{profileName})

	assert.Nil(t, requestErr, "Should load request correctly")
	if requestErr != nil {
		panic(requestErr)
	}

	assert.Equal(t, "some/path", request.URL, "Should load the correct URL")
}

func testErrorOnRequestNotFound(t *testing.T) {
	tempProfilesDir := SetupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "requests:\n  myRequest:\n    url: some/path\n"
	CreateTestProfile(profileName, profileContent, tempProfilesDir)

	_, requestErr := LoadRequestOptions("anotherRequest", []string{profileName})

	assert.NotNil(t, requestErr, "Should fail to load if request name not found")
}
