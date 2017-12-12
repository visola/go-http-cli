package daemon

import (
	"encoding/json"
	"net/http"

	"github.com/visola/go-http-cli/config"
	"github.com/visola/go-http-cli/ioutil"
)

// Handshake connects and sends a handshake request to the daemon. Return the version of the daemon
// that answered.
func Handshake() (int8, error) {
	var handshake HandshakeResponse

	if callDaemonError := callDaemon("/", "", &handshake); callDaemonError != nil {
		return 0, callDaemonError
	}

	return handshake.MajorVersion, nil
}

func callDaemon(path string, data string, unmarshalTo interface{}) error {
	method := "POST"

	if data == "" {
		method = "GET"
	}

	url := "http://localhost:" + string(DaemonPort) + path
	req, reqErr := http.NewRequest(method, url, nil)

	if reqErr != nil {
		return reqErr
	}

	if data != "" {
		req.Body = ioutil.CreateCloseableBufferString(data)
	}

	client := &http.Client{}
	response, responseErr := client.Do(req)

	if responseErr != nil {
		return responseErr
	}

	if unmarshalError := json.NewDecoder(response.Body).Decode(unmarshalTo); unmarshalError != nil {
		return unmarshalError
	}

	return nil
}
