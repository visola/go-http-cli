package session

import (
	"net/http"
)

// Session stores information for one session
type Session struct {
	Cookies   map[string]*http.Cookie
	Host      string
	Variables map[string]string
}
