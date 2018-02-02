package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/cli"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/output"
	"github.com/visola/go-http-cli/request"
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

	executionOptions := request.ExecutionOptions{
		FileToUpload:   options.FileToUpload,
		FollowLocation: options.FollowLocation,
		MaxRedirect:    options.MaxRedirect,
		ProfileNames:   options.Profiles,
		RequestName:    options.RequestName,
		Request: request.Request{
			Body:    options.Body,
			Headers: options.Headers,
			Method:  options.Method,
			URL:     options.URL,
		},
		Variables: options.Variables,
	}

	requestExecution, requestError := daemon.ExecuteRequest(executionOptions)

	if requestError != nil {
		color.Red("Error while executing request: %s", requestError)
		os.Exit(10)
	}

	for _, requestResponse := range requestExecution.RequestResponses {
		output.PrintRequest(requestResponse.Request)
		fmt.Println("")
		output.PrintResponse(requestResponse.Response)
	}

	if requestExecution.ErrorMessage != "" {
		color.Red("Error while executing request: %s", requestExecution.ErrorMessage)
		os.Exit(20)
	}
}
