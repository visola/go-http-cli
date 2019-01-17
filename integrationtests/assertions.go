package integrationtests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HasBody asserts that the request received the specified body
func HasBody(t *testing.T, req Request, body string) {
	assert.Equal(t, body, req.Body, "Should match body")
}

// HasMethod asserts that the request received the specified method
func HasMethod(t *testing.T, req Request, method string) {
	assert.Equal(t, method, req.Method, fmt.Sprintf("Method should be '%s'", method))
}

// HasPath asserts that the request received the specified path
func HasPath(t *testing.T, req Request, expectedPath string) {
	assert.Equal(t, expectedPath, req.Path, "Should match path")
}

// HasQueryParam checks if a request has the query parameter with the specified value
func HasQueryParam(t *testing.T, req Request, name string, value string) {
	paramValues := req.Query[name]
	assert.NotEmpty(t, paramValues, fmt.Sprintf("Expected query param: '%s'", name))
	if len(paramValues) > 0 {
		assert.Contains(t, paramValues, value, fmt.Sprintf("Param '%s' should include value '%s", name, value))
	}
}
