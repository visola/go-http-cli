package request

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/visola/go-http-cli/profile"
)

func TestConfigureRequest(t *testing.T) {
	t.Run("Simple request", testSimpleRequest)
	t.Run("Test with body", testWithBody)
	t.Run("Test with values", testWithValues)
	t.Run("Test POST with values", testPostWithValues)
	t.Run("Test with body and values", testWithBodyAndValues)
	t.Run("Test with profiles", testConfigureFromProfile)
}

func testSimpleRequest(t *testing.T) {
	req := Request{
		URL: "http://www.someserver.com/some/path",
	}

	configureRequest, err := ConfigureRequest(req)

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, req.URL, configureRequest.URL, "Should set passed in URL")
	assert.Equal(t, http.MethodGet, configureRequest.Method, "Should be set to GET")
}

func testWithBody(t *testing.T) {
	req := Request{
		Body: "Hello server!",
		URL:  "http://www.someserver.com/some/path",
	}

	configureRequest, err := ConfigureRequest(req)

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, req.URL, configureRequest.URL, "Should set passed in URL")
	assert.Equal(t, http.MethodPost, configureRequest.Method, "Should be set to POST")
	assert.Equal(t, req.Body, configureRequest.Body, "Should set body correctly")
}

func testPostWithValues(t *testing.T) {
	req := Request{
		Method: http.MethodPost,
		URL:    "http://www.someserver.com/some/path",
	}

	values := map[string][]string{
		"name": []string{"John"},
		"age":  []string{"20"},
	}

	configureRequest, err := ConfigureRequest(req, AddValues(values))

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, req.URL, configureRequest.URL, "Should set passed in URL")
	assert.Equal(t, req.Method, configureRequest.Method, "Should be set to method passed in")
	assert.Equal(t, `{"age":"20","name":"John"}`, configureRequest.Body, "Should set body as JSON")
}

func testWithValues(t *testing.T) {
	req := Request{
		URL: "http://www.someserver.com/some/path",
	}

	values := map[string][]string{
		"name": []string{"John"},
		"age":  []string{"20"},
	}

	configureRequest, err := ConfigureRequest(req, AddValues(values))

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, req.URL, configureRequest.URL, "Should set passed in URL")
	assert.Equal(t, http.MethodGet, configureRequest.Method, "Should be set to GET")
	assert.Equal(t, values, configureRequest.QueryParams, "Should set values as query parameters")
}

func testWithBodyAndValues(t *testing.T) {
	req := Request{
		Body: "Hello server!",
		URL:  "http://www.someserver.com/some/path",
	}

	values := map[string][]string{
		"age":  []string{"{age}"},
		"name": []string{"{name}"},
	}

	configureRequest, err := ConfigureRequest(req, AddValues(values))

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, req.URL, configureRequest.URL, "Should set passed in URL")
	assert.Equal(t, http.MethodPost, configureRequest.Method, "Should be set to POST")
	assert.Equal(t, req.Body, configureRequest.Body, "Should bet set to body")
	assert.Equal(t, values, configureRequest.QueryParams, "Should set query params")
}

func testConfigureFromProfile(t *testing.T) {
	profileName := "testProfile"
	profileContent := "baseURL: http://www.someserver.com/"

	profileContent = profileContent + "\nheaders:"
	profileContent = profileContent + "\n  Content-Type: application/json"
	profileContent = profileContent + "\n  Company-Id: '{companyId}'"
	profileContent = profileContent + "\n  X-Some-Header: '4321-4321-4321'" // This header will be overridden

	profileContent = profileContent + "\nvariables:"
	profileContent = profileContent + "\n  companyId: 1234"
	profileContent = profileContent + "\n  username: John Doe"

	profileContent = profileContent + "\nrequests:"
	profileContent = profileContent + "\n  withFile:"
	profileContent = profileContent + "\n    method: PUT"
	profileContent = profileContent + "\n    url: '/{companyId}/employee'"
	profileContent = profileContent + "\n    fileToUpload: test-body.yml"
	profileContent = profileContent + "\n    headers:"
	profileContent = profileContent + "\n      X-Some-Header: '1234-1234-1234'" // This will override the previously header

	testProfileDir := profile.SetupTestProfilesDir()
	profile.CreateTestProfile(profileName, profileContent, testProfileDir)

	jsonBody := `{"name":"John Doe","companyId":{companyId}}`
	profile.CreateTestProfile("test-body", jsonBody, testProfileDir)

	configureRequest, err := ConfigureRequest(Request{}, SetRequestName("withFile"), AddProfiles(profileName))

	assert.Nil(t, err, "Should not return an error")
	if err != nil {
		return
	}

	assert.Equal(t, "http://www.someserver.com/{companyId}/employee", configureRequest.URL, "Should build URL correctly")
	assert.Equal(t, 3, len(configureRequest.Headers), "Should configure all headers correctly")
	assert.Equal(t, http.MethodPut, configureRequest.Method, "Should set method from profile")
	assert.Equal(t, []string{"application/json"}, configureRequest.Headers["Content-Type"], "Should setup header from profile")
	assert.Equal(t, []string{"1234-1234-1234"}, configureRequest.Headers["X-Some-Header"], "Should override header correctly from request")
}
