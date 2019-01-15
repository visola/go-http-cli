package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/cli"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/output"
	"github.com/visola/go-http-cli/profile"
	"github.com/visola/go-http-cli/request"
)

func main() {
	ensureDaemon()

	options := parseCommandLineArguments()

	configureRequestOptions := request.CreateConfigureRequestOptions(
		request.AddProfiles(options.Profiles...),
		request.AddValues(options.Values),
		request.SetRequestName(options.RequestName),
	)

	mergedProfile, profileError := profile.LoadAndMergeProfiles(configureRequestOptions.ProfileNames)
	if profileError != nil {
		panic(profileError)
	}

	configuredRequest, configureError := request.ConfigureRequest(
		initializeRequest(options),
		&mergedProfile,
		configureRequestOptions,
	)

	if configureError != nil {
		panic(configureError)
	}

	executionContext := request.ExecutionContext{
		FollowLocation:  options.FollowLocation,
		MaxRedirect:     options.MaxRedirect,
		PostProcessCode: loadPostProcessScript(options, mergedProfile),
		ProfileNames:    options.Profiles,
		Request:         *configuredRequest,
		Variables:       options.Variables,
	}

	requestExecution, requestError := daemon.ExecuteRequest(executionContext)
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

func initializeRequest(options *cli.CommandLineOptions) request.Request {
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

	return unconfiguredRequest
}

func loadPostProcessScript(options *cli.CommandLineOptions, mergedProfiles profile.Options) request.PostProcessSourceCode {
	if options.PostProcessFile != "" {
		sourceCode, readErr := ioutil.ReadFile(options.PostProcessFile)
		if readErr != nil {
			panic(readErr)
		}
		return request.PostProcessSourceCode{
			SourceCode:     string(sourceCode),
			SourceFilePath: options.PostProcessFile,
		}
	}

	if options.RequestName != "" {
		namedRequest, findErr := profile.FindNamedRequest(&mergedProfiles, options.RequestName)
		if findErr != nil {
			panic(findErr)
		}

		if namedRequest.PostProcessScript != "" {
			return request.PostProcessSourceCode{
				SourceCode:     namedRequest.PostProcessScript,
				SourceFilePath: namedRequest.Source + ":requests." + options.RequestName + ".postProcessScript",
			}
		}
	}

	return request.PostProcessSourceCode{}
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
