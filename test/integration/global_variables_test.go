package integration

import (
	"net/http"
	"testing"
)

func TestGlobalVariables(t *testing.T) {
	t.Run("Can override global variables", WrapForIntegrationTest(testCanOverrideGlobalVariable))
	t.Run("Setting global variable works", WrapForIntegrationTest(testSetGlobalVariableWorks))
}

func testCanOverrideGlobalVariable(t *testing.T) {
	RunHTTP(t,
		"-V", "var1=value1",
	)

	RunHTTP(t,
		testServer.URL,
		"-X", "POST",
		"-V", "var1=value2",
		"name={var1}",
	)

	HasMethod(t, lastRequest, http.MethodPost)
	HasBody(t, lastRequest, `{"name":"value2"}`)
}

func testSetGlobalVariableWorks(t *testing.T) {
	RunHTTP(t,
		"-V", "var1=value1",
	)

	RunHTTP(t,
		testServer.URL,
		"-X", "POST",
		"name={var1}",
	)

	HasMethod(t, lastRequest, http.MethodPost)
	HasBody(t, lastRequest, `{"name":"value1"}`)
}
