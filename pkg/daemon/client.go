package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/visola/go-http-cli/pkg/ioutil"
	"github.com/visola/go-http-cli/pkg/request"
	"github.com/visola/go-http-cli/pkg/session"
)

// ExecuteRequest request the daemon to execute a request
func ExecuteRequest(executionContext request.ExecutionContext) (*RequestExecution, error) {
	dataAsBytes, marshalError := json.Marshal(executionContext)
	if marshalError != nil {
		return nil, marshalError
	}

	var requestExecution RequestExecution

	if callDaemonError := callDaemon("/request", string(dataAsBytes), &requestExecution); callDaemonError != nil {
		return nil, callDaemonError
	}

	return &requestExecution, nil
}

// Handshake connects and sends a handshake request to the daemon. Return the version of the daemon
// that answered.
func Handshake() (int8, error) {
	var handshake HandshakeResponse

	if callDaemonError := callDaemon("/", "", &handshake); callDaemonError != nil {
		return 0, callDaemonError
	}

	return handshake.MajorVersion, nil
}

// SetVariables sends variables to be set in the global session
func SetVariables(seVariablesRequest session.SetVariableRequest) error {
	dataAsBytes, marshalError := json.Marshal(seVariablesRequest)
	if marshalError != nil {
		return marshalError
	}

	if callDaemonError := callDaemon("/variables", string(dataAsBytes), nil); callDaemonError != nil {
		return callDaemonError
	}

	return nil
}

func callDaemon(path string, data string, unmarshalTo interface{}) error {
	method := http.MethodPost

	if data == "" {
		method = http.MethodGet
	}

	url := "http://localhost:" + string(DaemonPort) + path
	req, reqErr := http.NewRequest(method, url, nil)

	if reqErr != nil {
		return reqErr
	}

	if data != "" {
		req.Header.Add("Content-Type", "application/json")
		req.Body = ioutil.CreateCloseableBufferString(data)
	}

	client := &http.Client{}
	response, responseErr := client.Do(req)

	if responseErr != nil {
		return responseErr
	}

	if response.StatusCode != 200 {
		panic(fmt.Sprintf("Daemon responded with unexpected status code: %d - %s\nURL: %s, Method: %s", response.StatusCode, response.Status, url, method))
	}

	if unmarshalTo != nil {
		return json.NewDecoder(response.Body).Decode(unmarshalTo)
	}

	return nil
}
