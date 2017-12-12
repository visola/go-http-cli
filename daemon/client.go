package daemon

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/visola/go-http-cli/daemon/vo"
)

// Handshake connects and sends a handshake request to the daemon. Return the version of the daemon
// that answered.
func Handshake() (int8, error) {
	response, responseError := http.Get("http://localhost:" + string(DaemonPort))

	if responseError != nil {
		return 0, responseError
	}

	if response.StatusCode != 200 {
		return 0, errors.New("Daemon responded with unexpected status: " + string(response.StatusCode) + " - " + response.Status)
	}

	defer response.Body.Close()

	var handshake vo.HandshakeResponse

	if unmarshalError := json.NewDecoder(response.Body).Decode(&handshake); unmarshalError != nil {
		return 0, unmarshalError
	}

	return handshake.MajorVersion, nil
}
