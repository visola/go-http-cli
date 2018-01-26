package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProfile(t *testing.T) {
	t.Run("Loads basic authorization correctly", testLoadsBasicAuthorization)
}

func testLoadsBasicAuthorization(t *testing.T) {
	tempProfilesDir := SetupTestProfilesDir()

	profileName := "myProfile"
	profileContent := "auth:\n  type:  basic\n  username: myUsername\n  password: myPassword\n"
	CreateTestProfile(profileName, profileContent, tempProfilesDir)

	profile, requestErr := LoadProfile(profileName)

	assert.Nil(t, requestErr, "Should load request correctly")
	if requestErr != nil {
		panic(requestErr)
	}

	authValues, exists := profile.Headers["Authorization"]
	assert.True(t, exists, "Should have Authorization header")
	assert.Equal(t, 1, len(authValues), "Should have only one value")
}
