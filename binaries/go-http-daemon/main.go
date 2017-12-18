package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/options"
	"github.com/visola/go-http-cli/request"
)

var (
	log = logging.MustGetLogger("go-http-daemon")
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

	log.Fatal(server.Start(":" + string(daemon.DaemonPort)))
}

func configureLogging() {
	format := logging.MustStringFormatter(`%{color:bold}%{level} %{shortfunc} [%{time}]:%{color:reset} %{message}`)
	backend := logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format)
	logging.SetBackend(backend)
}

func executeRequest(c echo.Context) error {
	log.Debug("Execute request")

	requestOptions := new(options.RequestOptions)

	if parseRequestError := c.Bind(requestOptions); parseRequestError != nil {
		log.Error(parseRequestError)
		return parseRequestError
	}

	response, responseErr := request.ExecuteRequest(*requestOptions)

	if responseErr != nil {
		log.Error(responseErr)
		return responseErr
	}

	c.JSON(http.StatusOK, &daemon.ExecuteRequestResponse{
		RequestOptions: requestOptions,
		HTTPResponse:   response,
	})

	return nil
}

func handshake(c echo.Context) error {
	log.Debug("Handshake request")

	handshake := &daemon.HandshakeResponse{
		MajorVersion: daemon.DaemonMajorVersion,
		MinorVersion: daemon.DaemonMinorVersion,
	}

	c.JSON(http.StatusOK, handshake)
	return nil
}
