package integrationtests

import (
	"net/http"
	"testing"
)

func TestProfiles(t *testing.T) {
	t.Run("Profile With Named Request", WrapForIntegrationTest(testProfileWithNamedRequest))
}

func testProfileWithNamedRequest(t *testing.T) {
	CreateProfile("simple", `
baseURL: '{test-server}'

variables:
  companyId: 1234

requests:
  simple_request:
    url: /companies/{companyId}
`)

	RunHTTP(t, "+simple", "@simple_request")
	HasMethod(t, lastRequest, http.MethodGet)
	HasPath(t, lastRequest, "/companies/1234")
}
