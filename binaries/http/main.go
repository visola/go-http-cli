package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/options"
	"github.com/visola/go-http-cli/output"
)

func main() {
	if daemonErr := daemon.EnsureDaemon(); daemonErr != nil {
		panic(daemonErr)
	}

	options, err := options.ParseCommandLineOptions(os.Args[1:])

	if err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}

	executeRequestResponse, executeRequestError := daemon.ExecuteRequest(options)

	if executeRequestError != nil {
		color.Red("Error while executing request: %s", executeRequestError)
		os.Exit(10)
	}

	output.PrintRequest(executeRequestResponse.RequestOptions)
	output.PrintResponse(executeRequestResponse.HTTPResponse)
}
