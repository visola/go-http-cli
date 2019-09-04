package integration

import (
	"net/http"
	"testing"
)

func TestCookies(t *testing.T) {
	t.Run("Keeps track of cookies", WrapForIntegrationTest(testKeepsTrackOfCookies))
}

func testKeepsTrackOfCookies(t *testing.T) {
	replyWith := ReplyWith{
		Headers: map[string][]string{
			"Set-Cookie": []string{"someKey=someValue", "anotherCookie=someOtherValue"},
		},
	}

	prepareReply(replyWith)

	RunHTTP(t, testServer.URL)
	HasMethod(t, lastRequest, http.MethodGet)

	RunHTTP(t, testServer.URL)
	HasMethod(t, lastRequest, http.MethodGet)
	HasCookie(t, lastRequest, "someKey", "someValue")
	HasCookie(t, lastRequest, "anotherCookie", "someOtherValue")
}
