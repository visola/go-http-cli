package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/cli"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/output"
)

func main() {
	if daemonErr := daemon.EnsureDaemon(); daemonErr != nil {
		panic(daemonErr)
	}

	options, err := cli.ParseCommandLineOptions(os.Args[1:])

	if err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}

	requestResponses, requestError := daemon.ExecuteRequest(options)

	if requestError != nil {
		color.Red("Error while executing request: %s", requestError)
		os.Exit(10)
	}

	for _, requestResponse := range requestResponses {
		output.PrintRequest(requestResponse.Request)
		output.PrintResponse(requestResponse.Response)
	}
}
