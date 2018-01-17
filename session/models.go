package session

import (
	"net/http"
)

// Session stores information for one session
type Session struct {
	Jar http.CookieJar
}
