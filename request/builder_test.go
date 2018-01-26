package request

import (
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-http-cli/profile"
)

func TestBuildRequest(t *testing.T) {
	t.Run("Build request works", testRequestBuilder)
	t.Run("Build request with data from Profile", testRequestWithProfile)
}

func testRequestBuilder(t *testing.T) {
	request := Request{
		Body: `{"name":"{username},"companyId":{companyId}}`,
		Headers: map[string][]string{
			"Company-Id":   {"{companyId}"},
			"Content-Type": {"application/json"},
		},
		URL: "http://www.someserver.com/{companyId}/employee",
	}

	variables := map[string]string{
		"companyId": "1234",
		"username":  "John Doe",
	}

	httpReq, _, reqErr := BuildRequest(request, nil, variables)
	assert.Nil(t, reqErr, "Should create request")

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nCompany-Id: 1234\r\nContent-Type: application/json\r\n\r\n{\"name\":\"John Doe,\"companyId\":1234}",
		string(dump),
		"Should generate the expected dump",
	)
}

func testRequestWithProfile(t *testing.T) {
	profileName := "testProfile"
	profileContent := "baseURL: http://www.someserver.com/"
	profileContent = profileContent + "\n\nheaders:\n  Content-Type: application/json\n  Company-Id: '{companyId}'"
	profileContent = profileContent + "\n\nvariables:\n  companyId: 1234\n  username: John Doe"

	request := Request{
		Body: `{"name":"{username},"companyId":{companyId}}`,
		URL:  "/{companyId}/employee",
	}

	testProfileDir := profile.SetupTestProfilesDir()
	profile.CreateTestProfile(profileName, profileContent, testProfileDir)

	httpReq, _, reqErr := BuildRequest(request, []string{profileName}, nil)
	assert.Nil(t, reqErr, "Should create request")

	if reqErr != nil {
		panic(reqErr)
	}

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nCompany-Id: 1234\r\nContent-Type: application/json\r\n\r\n{\"name\":\"John Doe,\"companyId\":1234}",
		string(dump),
		"Should generate the expected dump",
	)
}
