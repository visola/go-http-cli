package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/request"
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

	server := echo.New()

	server.GET("/", handshake)
	server.POST("/request", executeRequest)

	log.Debugf("Daemon version %d.%d started and waiting for connections on port %s", daemon.DaemonMajorVersion, daemon.DaemonMinorVersion, daemon.DaemonPort)

	if writePIDError := daemon.WriteDaemonPID(); writePIDError != nil {
		panic(writePIDError)
	}

	go checkInteratcion()
	log.Fatal(server.Start(":" + string(daemon.DaemonPort)))
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

func executeRequest(c echo.Context) error {
	log.Debug("Execute request")
	lastInteraction = time.Now().UnixNano()

	requestExecution := daemon.RequestExecution{}
	daemonRequest := new(daemon.Request)

	if parseRequestError := c.Bind(daemonRequest); parseRequestError != nil {
		log.Error(parseRequestError)
		requestExecution.ErrorMessage = parseRequestError.Error()
		c.JSON(http.StatusOK, requestExecution)
		return nil
	}

	var req request.Request
	if daemonRequest.RequestName != "" {
		requestOptions, err := profile.LoadRequestOptions(daemonRequest.RequestName, daemonRequest.Profiles)

		if err != nil {
			requestExecution.ErrorMessage = err.Error()
			c.JSON(http.StatusOK, requestExecution)
			return nil
		}

		req = request.Request{
			Body:    requestOptions.Body,
			Headers: requestOptions.Headers,
			Method:  requestOptions.Method,
			URL:     requestOptions.URL,
		}

		req.Merge(daemonRequest.ToRequest())
	} else {
		req = daemonRequest.ToRequest()
	}

	requestResponses, responseErr := request.ExecuteRequest(req, daemonRequest.Profiles, daemonRequest.Variables)
	requestExecution.RequestResponses = requestResponses

	if responseErr != nil {
		log.Error(responseErr)
		requestExecution.ErrorMessage = responseErr.Error()
	}

	c.JSON(http.StatusOK, requestExecution)
	return nil
}

func handshake(c echo.Context) error {
	log.Debug("Handshake request")
	lastInteraction = time.Now().UnixNano()

	handshake := &daemon.HandshakeResponse{
		MajorVersion: daemon.DaemonMajorVersion,
		MinorVersion: daemon.DaemonMinorVersion,
	}

	c.JSON(http.StatusOK, handshake)
	return nil
}
