package request

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/visola/go-http-cli/profile"
)

func TestConfigureRequest(t *testing.T) {
	profileName := "testProfile"
	profileContent := "baseURL: http://www.someserver.com/"

	profileContent = profileContent + "\nheaders:"
	profileContent = profileContent + "\n  Content-Type: application/json"
	profileContent = profileContent + "\n  Company-Id: '{companyId}'"
	profileContent = profileContent + "\n  X-Some-Header: '4321-4321-4321'" // This header will be overriden

	profileContent = profileContent + "\nvariables:"
	profileContent = profileContent + "\n  companyId: 1234"
	profileContent = profileContent + "\n  username: John Doe"

	profileContent = profileContent + "\nrequests:"
	profileContent = profileContent + "\n  withFile:"
	profileContent = profileContent + "\n    url: '/{companyId}/employee'"
	profileContent = profileContent + "\n    fileToUpload: test-body.yml"
	profileContent = profileContent + "\n    headers:"
	profileContent = profileContent + "\n      X-Some-Header: '1234-1234-1234'" // This will override the previously header

	testProfileDir := profile.SetupTestProfilesDir()
	profile.CreateTestProfile(profileName, profileContent, testProfileDir)

	jsonBody := `{"name":"John Doe","companyId":{companyId}}`
	profile.CreateTestProfile("test-body", jsonBody, testProfileDir)

	configureRequest, err := ConfigureRequest(Request{}, "withFile", ExecutionOptions{ProfileNames: []string{profileName}})

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, "http://www.someserver.com/{companyId}/employee", configureRequest.URL, "Should build URL correctly")

	assert.Equal(t, 3, len(configureRequest.Headers), "Should configure all headers correctly")
	assert.Equal(t, []string{"application/json"}, configureRequest.Headers["Content-Type"], "Should setup header from profile")
	assert.Equal(t, []string{"1234-1234-1234"}, configureRequest.Headers["X-Some-Header"], "Should override header correctly from request")
}
