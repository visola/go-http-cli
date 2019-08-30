package session

import (
	"net/http"
)

// Session stores information for one session
type Session struct {
	Cookies   map[string]*http.Cookie
	Variables map[string]string
}
