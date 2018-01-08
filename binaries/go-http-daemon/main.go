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

	daemonRequest := new(daemon.Request)

	if parseRequestError := c.Bind(daemonRequest); parseRequestError != nil {
		log.Error(parseRequestError)
		return parseRequestError
	}

	var req request.Request
	if daemonRequest.RequestName != "" {
		requestOptions, err := profile.LoadRequestOptions(daemonRequest.RequestName, daemonRequest.Profiles)

		if err != nil {
			return err
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

	requestResponsePair, responseErr := request.ExecuteRequest(req, daemonRequest.Profiles, daemonRequest.Variables)

	if responseErr != nil {
		log.Error(responseErr)
		return responseErr
	}

	c.JSON(http.StatusOK, requestResponsePair)
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
