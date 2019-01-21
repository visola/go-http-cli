package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
	"github.com/visola/go-http-cli/pkg/daemon"
	"github.com/visola/go-http-cli/pkg/request"
)

var (
	log             = logging.MustGetLogger("go-http-daemon")
	lastInteraction = time.Now().UnixNano()
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--kill" {
		log.Info("Killing daemon...")
		daemon.KillDaemon()
		return
	}

	configureLogging()

	server := mux.NewRouter()
	server.HandleFunc("/", handshake).Methods(http.MethodGet)
	server.HandleFunc("/request", executeRequest).Methods(http.MethodPost)

	log.Debugf("Daemon version %d.%d started and waiting for connections on port %s", daemon.DaemonMajorVersion, daemon.DaemonMinorVersion, daemon.DaemonPort)

	if writePIDError := daemon.WriteDaemonPID(); writePIDError != nil {
		panic(writePIDError)
	}

	go checkInteratcion()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", daemon.DaemonPort), server))
}

func configureLogging() {
	format := logging.MustStringFormatter(`%{color:bold}%{level} %{shortfunc} [%{time}]:%{color:reset} %{message}`)
	backend := logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format)
	logging.SetBackend(backend)
}

func checkInteratcion() {
	for {
		now := time.Now().UnixNano()
		if now-lastInteraction > (30 * time.Minute).Nanoseconds() {
			log.Info("Too quiet around here, shutting down.")
			os.Exit(0)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func executeRequest(w http.ResponseWriter, req *http.Request) {
	log.Debug("Execute request")
	lastInteraction = time.Now().UnixNano()

	requestExecution := daemon.RequestExecution{}
	var executionContext request.ExecutionContext

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	if parseRequestError := decoder.Decode(&executionContext); parseRequestError != nil {
		log.Error(parseRequestError)
		requestExecution.ErrorMessage = parseRequestError.Error()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(requestExecution)
		return
	}

	requestResponses, responseErr := request.ExecuteRequestLoop(executionContext)
	requestExecution.RequestResponses = requestResponses

	if responseErr != nil {
		log.Error(responseErr)
		requestExecution.ErrorMessage = responseErr.Error()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requestExecution)
}

func handshake(w http.ResponseWriter, req *http.Request) {
	log.Debug("Handshake request")
	lastInteraction = time.Now().UnixNano()

	handshake := &daemon.HandshakeResponse{
		MajorVersion: daemon.DaemonMajorVersion,
		MinorVersion: daemon.DaemonMinorVersion,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handshake)
}
