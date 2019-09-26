package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/pkg/cli"
	"github.com/visola/go-http-cli/pkg/daemon"
	"github.com/visola/go-http-cli/pkg/model"
	"github.com/visola/go-http-cli/pkg/output"
	"github.com/visola/go-http-cli/pkg/profile"
	"github.com/visola/go-http-cli/pkg/request"
	"github.com/visola/go-http-cli/pkg/session"
)

func main() {
	ensureDaemon()

	options := parseCommandLineArguments()

	checkForSetVariableRequest(options)

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

	configuredRequest.PostProcessCode = loadPostProcessScript(options, mergedProfile)

	executionContext := request.ExecutionContext{
		FollowLocation: options.FollowLocation,
		MaxRedirect:    options.MaxRedirect,
		ProfileNames:   options.Profiles,
		Request:        *configuredRequest,
		Variables:      options.Variables,
	}

	requestExecution, requestError := daemon.ExecuteRequest(executionContext)
	if requestError != nil {
		color.Red("Error while executing request: %s", requestError)
		os.Exit(10)
	}

	printOutput(requestExecution, options)
}

func checkForSetVariableRequest(options *cli.CommandLineOptions) {
	// If passed variables, no profiles and no URL, then it sets a variable to global session
	if len(options.Variables) > 0 && len(options.Profiles) == 0 && options.URL == "" {
		setVariableRequest := session.SetVariableRequest{
			Values: make([]model.KeyValuePair, 0),
		}

		for name, value := range options.Variables {
			setVariableRequest.Values = append(setVariableRequest.Values, model.KeyValuePair{
				Name:  name,
				Value: value,
			})
		}

		if err := daemon.SetVariables(setVariableRequest); err != nil {
			panic(err)
		}

		color.Green("Variables set.")
		os.Exit(0)
	}
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
	failedRequest := 0

	for _, requestResponse := range requestExecution.RequestResponses {
		output.PrintRequest(requestResponse.Request)
		fmt.Println("")
		output.PrintResponse(requestResponse.Response)
		fmt.Println("")

		if requestResponse.Response.StatusCode >= http.StatusBadRequest {
			failedRequest++
		}

		if options.OutputFile != "" && requestResponse.Response.Body != "" {
			outWriteErr := ioutil.WriteFile(options.OutputFile, []byte(requestResponse.Response.Body), 0644)
			if outWriteErr != nil {
				color.Red("Error while writing to output file: %s", outWriteErr)
			}
		}

		if requestResponse.PostProcessOutput != "" {
			postProcessColor := color.New(color.FgBlue).PrintfFunc()
			postProcessColor("\n-- Post processing output --")
			postProcessColor("\n%s", requestResponse.PostProcessOutput)
			postProcessColor("\n-- End of output --\n")
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

	if len(requestExecution.RequestResponses) > 1 {
		color.Green("Number of requests: %d\n", len(requestExecution.RequestResponses))
		if failedRequest > 0 {
			color.Red("Number of failed requests: %d\n", failedRequest)
		}
	}

	os.Exit(exitCode)
}
