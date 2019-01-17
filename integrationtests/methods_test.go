package integrationtests

import (
	"net/http"
	"testing"
)

func TestMethods(t *testing.T) {
	t.Run("Method GET", WrapWithKillDamon(WrapWithTestServer(testGet)))
}

func testGet(t *testing.T) {
	userID := "1234"
	companyID := "7890"
	timestamp := "20181201083035"
	path := "/users/{userID}?timestamp={timestamp}"

	RunHTTP(
		t,
		"-V", "companyID="+companyID,
		"-V", "userID="+userID,
		"-V", "timestamp="+timestamp,
		testServer.URL+path,
		"companyID={companyID}",
	)

	HasMethod(t, lastRequest, http.MethodGet)
	HasPath(t, lastRequest, "/users/"+userID)
	HasQueryParam(t, lastRequest, "companyID", companyID)
	HasQueryParam(t, lastRequest, "timestamp", timestamp)
}
