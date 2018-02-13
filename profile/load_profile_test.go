package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadProfile(t *testing.T) {
	tempProfilesDir := SetupTestProfilesDir()

	// Profile adds content type header
	jsonProfileName := "json"
	jsonProfileContent := "headers:\n  Content-Type: application/json"
	CreateTestProfile(jsonProfileName, jsonProfileContent, tempProfilesDir)

	// Profile adds auth
	authProfileName := "auth"
	authProfileContent := "auth:\n  type:  basic\n  username: myUsername\n  password: myPassword\n"
	CreateTestProfile(authProfileName, authProfileContent, tempProfilesDir)

	// Profile add more info and import the two above
	profileName := "profile"
	profileContent := "import:\n  - json\n  - auth"
	profileContent = profileContent + "\n\nbaseURL: http://www.someserver.com"
	CreateTestProfile(profileName, profileContent, tempProfilesDir)

	profile, requestErr := LoadProfile(profileName)

	assert.Nil(t, requestErr, "Should load request correctly")
	if requestErr != nil {
		panic(requestErr)
	}

	assert.Equal(t, "http://www.someserver.com", profile.BaseURL, "Should set the base URL correctly")

	authValues, exists := profile.Headers["Authorization"]
	assert.True(t, exists, "Should have Authorization header")
	assert.Equal(t, 1, len(authValues), "Should have only one value")
	assert.Equal(t, "Basic bXlVc2VybmFtZTpteVBhc3N3b3Jk", authValues[0], "Should set the auth value correctly")

	contentTypeValue, exists := profile.Headers["Content-Type"]
	assert.True(t, exists, "Should have Content-Type header")
	assert.Equal(t, 1, len(contentTypeValue), "Should have only one value")
	assert.Equal(t, "application/json", contentTypeValue[0], "Should set the value correctly")
}
