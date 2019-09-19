package integration

import (
	"net/http"
	"os"
	"testing"
)

func TestPostProcess(t *testing.T) {
	t.Run("Add request as object", WrapForIntegrationTest(testAddRequestAsObjectFromScript))
	t.Run("Add request by name", WrapForIntegrationTest(testAddRequestByNameFromScript))
	t.Run("Add request by URL", WrapForIntegrationTest(testAddRequestByURLFromScript))
	t.Run("Set variable", WrapForIntegrationTest(testSetVariableFromScript))
}

func testAddRequestAsObjectFromScript(t *testing.T) {
	CreateProfile("test", `
baseURL: '{test-server}'

variables:
  companyId: 1234
`)

	postProcessScript := `
		addRequest({
			Body: '{"name":"Some Company"}',
			Method: 'POST',
			URL: '/companies/{companyId}',
		});
	`

	WithTempFile(t, postProcessScript, func (tempFile *os.File) {
		RunHTTP(t, "+test", "--post-process", tempFile.Name(), testServer.URL+"/hello")

		HasRequestCount(t, 2)
		HasBody(t, allRequests[1], `{"name":"Some Company"}`)
		HasMethod(t, allRequests[1], http.MethodPost)
		HasPath(t, allRequests[1], "/companies/1234")
	})
}

func testAddRequestByNameFromScript(t *testing.T) {
	CreateProfile("test", `
baseURL: '{test-server}'

requests:
  someRequest:
    body: 'Some Request'
    method: POST
    url: /some/place
`)

	postProcessScript := "addRequest('@someRequest');"
	WithTempFile(t, postProcessScript, func (tempFile *os.File) {
		RunHTTP(t, "+test", "--post-process", tempFile.Name(), testServer.URL+"/hello")

		HasRequestCount(t, 2)
		HasBody(t, allRequests[1], "Some Request")
		HasMethod(t, allRequests[1], http.MethodPost)
		HasPath(t, allRequests[1], "/some/place")
	})
}

func testAddRequestByURLFromScript(t *testing.T) {
	CreateProfile("test", `
baseURL: '{test-server}'

variables:
  companyId: 1234
`)

	postProcessScript := "addRequest('/companies/{companyId}');"
	WithTempFile(t, postProcessScript, func (tempFile *os.File) {
		RunHTTP(t, "+test", "--post-process", tempFile.Name(), testServer.URL+"/hello")

		HasRequestCount(t, 2)
		HasPath(t, allRequests[1], "/companies/1234")
	})
}

func testSetVariableFromScript(t *testing.T) {
	// Post process will get token from body, after parsing JSON
	postProcessScript := `
		var body = JSON.parse(response.Body);
		addVariable('token', body.token);
	`
	
	WithTempFile(t, postProcessScript, func (tempFile *os.File) {
		// Set reply body
		token := "not-real-token"
		replyWith := ReplyWith{
			Body: `{"token":"`+token+`"}`,
		}
		prepareReply(replyWith)
	
		// Run "login" and post process response
		RunHTTP(t, "-X", "POST", "--post-process", tempFile.Name(), testServer.URL+"/login")
	
		// Second request uses the token from post process
		RunHTTP(t, "-H", "Authorization: Bearer {token}", testServer.URL+"/dashboard")
	
		// Check that the header was set to the correct value
		HasHeader(t, lastRequest, "Authorization", "Bearer "+token)
	})
}
