package session

import (
	"net/http"
	"sync"
)

var sessions = map[string]*Session{
	// Global session
	"": &Session{
		Cookies:   make(map[string]*http.Cookie, 0),
		Variables: make(map[string]string),
	},
}

var sessionMutex = &sync.Mutex{}

// Get a session based on a URL.
func Get(host string) *Session {
	session := ensureSession(host)

	if host != "" {
		mergedSession := mergeSessions(sessions[""], session)
		mergedSession.Host = host
		return mergedSession
	}

	return session
}

// SetCookie sets a cookie into a session, creating one if it doesn't exist
func SetCookie(host string, cookie *http.Cookie) {
	ensureSession(host).Cookies[cookie.Name] = cookie
}

// SetGlobalVariable sets a value that will be merged into all sessions
func SetGlobalVariable(name, value string) {
	sessions[""].Variables[name] = value
}

// SetVariable sets a value to a specific variable, creating a session if it doesn't exists
func SetVariable(host, name, value string) {
	ensureSession(host).Variables[name] = value
}

func ensureSession(host string) *Session {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	session, exists := sessions[host]
	if !exists {
		session = &Session{
			Cookies:   make(map[string]*http.Cookie, 0),
			Host:      host,
			Variables: make(map[string]string),
		}

		sessions[host] = session
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
