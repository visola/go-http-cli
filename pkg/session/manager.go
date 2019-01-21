package session

import (
	"net/http"
	"sync"
)

var managerInstance = manager{
	sessions: make(map[string]*Session),
}

var sessionMutex = &sync.Mutex{}

// manager holds sessions based on domain.
type manager struct {
	sessions map[string]*Session
}

// Get a session based on a URL.
func Get(host string) (*Session, error) {
	sessionMutex.Lock()
	session, exists := managerInstance.sessions[host]

	if !exists {
		session = &Session{
			Cookies:   make([]*http.Cookie, 0),
			Variables: make(map[string]string),
		}

		managerInstance.sessions[host] = session
	}

	sessionMutex.Unlock()
	return session, nil
}
