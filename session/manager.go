package session

import (
	"net/http/cookiejar"
	"net/url"
)

var managerInstance = manager{
	sessions: make(map[string]Session),
}

// manager holds sessions based on domain.
type manager struct {
	sessions map[string]Session
}

// Get a session based on a URL.
func Get(url url.URL) (*Session, error) {
	host := url.Hostname()
	session, exists := managerInstance.sessions[host]

	if !exists {
		cookieJar, jarError := cookiejar.New(nil)
		if jarError != nil {
			return nil, jarError
		}

		session = Session{
			Jar: cookieJar,
		}

		managerInstance.sessions[host] = session
	}

	return &session, nil
}
