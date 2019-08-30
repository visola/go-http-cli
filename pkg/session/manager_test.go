package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreserveCookie(t *testing.T) {
	sessionName := "test"
	session := Get(sessionName)

	cookie := &http.Cookie{
		Name:  "Delicious",
		Value: "Yes!",
	}

	session.Cookies["Delicious"] = cookie

	session2 := Get(sessionName)
	cookie2, exists := session2.Cookies[cookie.Name]
	assert.True(t, exists, "Should keep cookie")
	if exists {
		assert.Equal(t, cookie2.Value, cookie.Value, "Should contain the same values")
	}
}
