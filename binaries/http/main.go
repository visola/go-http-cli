package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/cli"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/output"
	"github.com/visola/go-http-cli/request"
)

func main() {
	ensureDaemon()

	options := parseCommandLineArguments()
	configuredRequest := initializeAndConfigureRequest(options)

	executionOptions := request.ExecutionOptions{
		FollowLocation:  options.FollowLocation,
		MaxRedirect:     options.MaxRedirect,
		PostProcessFile: options.PostProcessFile,
		ProfileNames:    options.Profiles,
		Request:         *configuredRequest,
		Variables:       options.Variables,
	}

	requestExecution, requestError := daemon.ExecuteRequest(executionOptions)
	if requestError != nil {
		color.Red("Error while executing request: %s", requestError)
		os.Exit(10)
	}

	printOutput(requestExecution, options)
}

func ensureDaemon() {
	if daemonErr := daemon.EnsureDaemon(); daemonErr != nil {
		panic(daemonErr)
	}
}

func initializeAndConfigureRequest(options *cli.CommandLineOptions) *request.Request {
	unconfiguredRequest := request.Request{
		Body:    options.Body,
		Headers: options.Headers,
		Method:  options.Method,
		URL:     options.URL,
	}

	loadBodyError := unconfiguredRequest.LoadBodyFromFile(options.FileToUpload)
	if loadBodyError != nil {
		panic(loadBodyError)
	}

	configuredRequest, configureError := request.ConfigureRequest(
		unconfiguredRequest,
		request.AddProfiles(options.Profiles...),
		request.AddValues(options.Values),
		request.SetRequestName(options.RequestName),
	)

	if configureError != nil {
		panic(configureError)
	}

	return configuredRequest
}

func parseCommandLineArguments() *cli.CommandLineOptions {
	options, err := cli.ParseCommandLineOptions(os.Args[1:])
	if err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}
	return options
}

func printOutput(requestExecution *daemon.RequestExecution, options *cli.CommandLineOptions) {
	exitCode := 0
	for _, requestResponse := range requestExecution.RequestResponses {
		output.PrintRequest(requestResponse.Request)
		fmt.Println("")
		output.PrintResponse(requestResponse.Response)
		if options.OutputFile != "" && requestResponse.Response.Body != "" {
			outWriteErr := ioutil.WriteFile(options.OutputFile, []byte(requestResponse.Response.Body), 0644)
			if outWriteErr != nil {
				color.Red("Error while writing to output file: %s", outWriteErr)
			}
		}

		if requestResponse.PostProcessOutput != "" {
			fmt.Println("\n -- Post processing output --")
			fmt.Println(requestResponse.PostProcessOutput)
		}

		if requestResponse.PostProcessError != "" {
			color.Red("Error post processing request: %s", requestResponse.PostProcessError)
			exitCode = 30
		}
	}

	if requestExecution.ErrorMessage != "" {
		color.Red("Error while executing request: %s", requestExecution.ErrorMessage)
		exitCode = 20
	}

	os.Exit(exitCode)
}
