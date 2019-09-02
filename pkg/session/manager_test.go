package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionManagement(t *testing.T) {
	t.Run("Sets cookie correctly", testSetsCookie)
	t.Run("Sets value correctly", testSetsValue)
	t.Run("Sets value globally", testSetsValueGlobally)
	t.Run("Sets value overrides global", testVariableOverridesGlobal)
}

func testSetsCookie(t *testing.T) {
	sessionName := "test"

	cookie := &http.Cookie{
		Name:  "Delicious",
		Value: "Yes!",
	}

	SetCookie(sessionName, cookie)

	session := Get(sessionName)

	assert.Equal(t, sessionName, session.Host, "Should have correct session name")

	cookie2, exists := session.Cookies[cookie.Name]
	assert.True(t, exists, "Should set cookie")
	if exists {
		assert.Equal(t, cookie2.Value, cookie.Value, "Should contain the same values")
	}
}

func testSetsValue(t *testing.T) {
	sessionName := "test"
	variableName := "variable"
	variableValue := "value"
	SetVariable(sessionName, variableName, variableValue)

	session := Get(sessionName)
	value, exists := session.Variables[variableName]
	assert.True(t, exists, "Should set variable")
	if exists {
		assert.Equal(t, variableValue, value, "Should have the correct value")
	}
}

func testSetsValueGlobally(t *testing.T) {
	sessionName := "test"
	variableName := "variable"
	variableValue := "value"
	SetGlobalVariable(variableName, variableValue)

	session := Get(sessionName)
	value, exists := session.Variables[variableName]
	assert.True(t, exists, "Should set variable to all sessions")
	if exists {
		assert.Equal(t, variableValue, value, "Should have the correct value")
	}
}

func testVariableOverridesGlobal(t *testing.T) {
	sessionName := "test"
	variableName := "variable"
	variableValue := "value"

	SetGlobalVariable(variableName, "Global Value")
	SetVariable(sessionName, variableName, variableValue)

	session := Get(sessionName)
	value, exists := session.Variables[variableName]
	assert.True(t, exists, "Should set variable to all sessions")
	if exists {
		assert.Equal(t, variableValue, value, "Should have the correct value")
	}
}
