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

	executedRequestResponse, executeRequestError := daemon.ExecuteRequest(options)

	if executeRequestError != nil {
		color.Red("Error while executing request: %s", executeRequestError)
		os.Exit(10)
	}

	output.PrintRequest(executedRequestResponse.Request)
	output.PrintResponse(executedRequestResponse.Response)
}
