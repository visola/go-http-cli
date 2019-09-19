package integration

import (
	"os"
	"testing"
)

func TestPostProcess(t *testing.T) {
	t.Run("Can set variable", WrapForIntegrationTest(testCanSetVariableFromScript))
}

func testCanSetVariableFromScript(t *testing.T) {
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
