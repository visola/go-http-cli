package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/op/go-logging"

	"github.com/visola/go-http-cli/daemon"
)

var (
	log = logging.MustGetLogger("go-http-daemon")
)

func main() {
	configureLogging()

	http.HandleFunc("/", handshake)

	log.Debugf("Daemon version %d.%d started and waiting for connections on port %s", daemon.DaemonMajorVersion, daemon.DaemonMinorVersion, daemon.DaemonPort)

	if writePIDError := daemon.WriteDaemonPID(); writePIDError != nil {
		panic(writePIDError)
	}

	log.Fatal(http.ListenAndServe(":"+string(daemon.DaemonPort), nil))
}

func configureLogging() {
	format := logging.MustStringFormatter(`%{color:bold}%{level:.4s} %{shortfunc} [%{time}]:%{color:reset} %{message}`)
	backend := logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format)
	logging.SetBackend(backend)
}

func handshake(response http.ResponseWriter, request *http.Request) {
	log.Debug("Handshake request")

	handshake := &daemon.HandshakeResponse{
		MajorVersion: daemon.DaemonMajorVersion,
		MinorVersion: daemon.DaemonMinorVersion,
	}
	json.NewEncoder(response).Encode(handshake)
}

