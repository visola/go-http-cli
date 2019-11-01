package integration

import (
	"net/http"
	"testing"
)

func TestProfiles(t *testing.T) {
	t.Run("Profiles with inheritance", WrapForIntegrationTest(testProfileInheritance))
	t.Run("Profile With Named Request", WrapForIntegrationTest(testProfileWithNamedRequest))
	t.Run("Profile with POST using form and variables", WrapForIntegrationTest(testProfileWithVariableInForm))
	t.Run("Profile with POST using string form and variables", WrapForIntegrationTest(testProfileWithVariableInStringForm))
}

func testProfileInheritance(t *testing.T) {
	CreateProfile("parent_profile", `
baseURL: https://doesnt.matter.com/

requests:
  post_with_body:
    body: '{"name":"John Doe"}'
`)

	CreateProfile("child_profile", `
baseURL: '{test-server}'
import:
  - parent_profile
`)

	RunHTTP(t, "+child_profile", "@post_with_body")

	HasBody(t, lastRequest, `{"name":"John Doe"}`)
	HasHeader(t, lastRequest, "Content-Type", "application/json")
	HasMethod(t, lastRequest, http.MethodPost)
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

func testProfileWithVariableInForm(t *testing.T) {
	CreateProfile("formWithVariable", `
baseURL: '{test-server}'

variables:
  companyId: 1234

requests:
  auth:
    url: /auth
    method: POST
    headers:
      Content-Type: application/x-www-form-urlencoded
    values:
      username: '{username}'
      password: '{password}'
`)

	username := "johndoe@gmail.com"
	password := "FK5wHc2!&i"

	RunHTTP(t, "+formWithVariable", "@auth", "-V", "username="+username, "-V", "password="+password)
	HasMethod(t, lastRequest, http.MethodPost)
	HasPath(t, lastRequest, "/auth")
	HasBody(t, lastRequest, "password=FK5wHc2%21%26i&username=johndoe%40gmail.com")
}

func testProfileWithVariableInStringForm(t *testing.T) {
	CreateProfile("formWithVariable", `
baseURL: '{test-server}'

variables:
  companyId: 1234

requests:
  auth:
    url: /auth
    method: POST
    headers:
      Content-Type: application/x-www-form-urlencoded
    body:
      username={username}&password={password}
`)

	username := "johndoe@gmail.com"
	password := "FK5wHc2!&i"

	RunHTTP(t, "+formWithVariable", "@auth", "-V", "username="+username, "-V", "password="+password)
	HasMethod(t, lastRequest, http.MethodPost)
	HasPath(t, lastRequest, "/auth")
	HasBody(t, lastRequest, "password=FK5wHc2%21%26i&username=johndoe%40gmail.com")
}
