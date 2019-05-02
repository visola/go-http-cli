package integration

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
	value2 := "Some Value for a Variable"
	path := "/users/{userID}?timestamp={timestamp}"

	RunHTTP(
		t,
		"-H", "Authorization: Basic QXJlIHlvdSBraWRkaW5nIG1lPw==",
		"-H", "Multi-Value=value1",
		"-H", "Multi-Value={value2}",
		"-V", "companyID="+companyID,
		"-V", "userID="+userID,
		"-V", "timestamp="+timestamp,
		"-V", "value2="+value2,
		"-X", method,
		testServer.URL+path,
		"createdAfter=2018-12-01 10:00:00+0400",
		"companyID={companyID}",
	)

	HasMethod(t, lastRequest, method)
	HasHeader(t, lastRequest, "Authorization", "Basic QXJlIHlvdSBraWRkaW5nIG1lPw==")
	HasHeader(t, lastRequest, "Multi-Value", "value1")
	HasHeader(t, lastRequest, "Multi-Value", value2)
	HasPath(t, lastRequest, "/users/"+userID)
	HasQueryParam(t, lastRequest, "timestamp", timestamp)

	if method == http.MethodGet {
		HasQueryParam(t, lastRequest, "companyID", companyID)
		HasQueryParam(t, lastRequest, "createdAfter", "2018-12-01 10:00:00+0400")
	} else {
		HasBody(t, lastRequest, `{"companyID":"`+companyID+`","createdAfter":"2018-12-01 10:00:00+0400"}`)
	}
}
