package integrationtests

import (
	"net/http"
	"testing"
)

func TestMethods(t *testing.T) {
	methods := []string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}

	for _, method := range methods {
		t.Run("Test "+method, WrapForIntegrationTest(buildTestFunc(method)))
	}
}

func buildTestFunc(method string) func(*testing.T) {
	return func(t *testing.T) {
		testMethod(t, method)
	}
}

func testMethod(t *testing.T, method string) {
	userID := "1234"
	companyID := "7890"
	timestamp := "20181201083035"
	path := "/users/{userID}?timestamp={timestamp}"

	RunHTTP(
		t,
		"-V", "companyID="+companyID,
		"-V", "userID="+userID,
		"-V", "timestamp="+timestamp,
		"-X", method,
		testServer.URL+path,
		"companyID={companyID}",
	)

	HasMethod(t, lastRequest, method)
	HasPath(t, lastRequest, "/users/"+userID)
	HasQueryParam(t, lastRequest, "timestamp", timestamp)

	if method == http.MethodGet {
		HasQueryParam(t, lastRequest, "companyID", companyID)
	} else {
		HasBody(t, lastRequest, `{"companyID":"`+companyID+`"}`)
	}
}
