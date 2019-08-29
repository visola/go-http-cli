package session

import (
	"net/http"
	"sync"
)

var sessions = make(map[string]*Session)
var sessionMutex = &sync.Mutex{}

// Get a session based on a URL.
func Get(host string) *Session {
	sessionMutex.Lock()
	session, exists := sessions[host]

	if !exists {
		session = &Session{
			Cookies:   make([]*http.Cookie, 0),
			Variables: make(map[string]string),
		}

		sessions[host] = session
	}

	sessionMutex.Unlock()
	return session
}
