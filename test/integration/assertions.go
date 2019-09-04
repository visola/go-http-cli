package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HasBody asserts that the request received the specified body
func HasBody(t *testing.T, req Request, body string) {
	assert.Equal(t, body, req.Body, "Should match body")
}

// HasCookie asserts that the request received the specific cookies
func HasCookie(t *testing.T, req Request, name string, value string) {
	assert.True(t, len(req.Cookies) > 0, "Request does not contain any cookies.")
	if len(req.Cookies) > 0 {
		for _, cookie := range req.Cookies {
			if cookie.Name == name && cookie.Value == value {
				return
			}
		}
		assert.Fail(t, fmt.Sprintf("Cookies do not contain %s=%s => %s", name, value, req.Cookies))
	}
}

// HasHeader asserts that the rquest received the specified header with value
func HasHeader(t *testing.T, req Request, name string, value string) {
	checkMapOfArrayOfStrings(t, req.Headers, name, value, "header")
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
	checkMapOfArrayOfStrings(t, req.Query, name, value, "query param")
}

func checkMapOfArrayOfStrings(t *testing.T, toCheck map[string][]string, name string, value string, alias string) {
	values, exists := toCheck[name]
	if !exists {
		values = toCheck[strings.ToLower(name)]
	}

	assert.NotEmpty(t, values, fmt.Sprintf("Expected %s: '%s'", alias, name))
	if len(values) > 0 {
		containsCaseInsensitive(t, values, value, fmt.Sprintf("%s '%s' should include value '%s", alias, name, value))
	}
}

func containsCaseInsensitive(t *testing.T, values []string, expectedValue, message string) {
	expectedValue = strings.ToLower(expectedValue)
	for _, value := range values {
		if strings.ToLower(value) == expectedValue {
			return
		}
	}
	assert.Fail(t, message)
}
