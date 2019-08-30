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
			Cookies:   make(map[string]*http.Cookie, 0),
			Variables: make(map[string]string),
		}

		sessions[host] = session
	}

	sessionMutex.Unlock()

	if host != "" {
		session = mergeSessions(session, Get(""))
	}

	return session
}

func mergeSessions(sessions ...*Session) *Session {
	finalSession := &Session{
		Cookies:   make(map[string]*http.Cookie, 0),
		Variables: make(map[string]string),
	}

	for _, session := range sessions {
		for n, v := range session.Cookies {
			finalSession.Cookies[n] = v
		}
		for n, v := range session.Variables {
			finalSession.Variables[n] = v
		}
	}

	return finalSession
}
