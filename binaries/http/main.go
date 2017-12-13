package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/daemon"
	"github.com/visola/go-http-cli/options"
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

	color.Green("Response status code: %d", executeRequestResponse.StatusCode)

	/*
		req, reqErr := request.BuildRequest(configuration)
		if reqErr != nil {
			color.Red("Error while creating request: %s", reqErr)
			os.Exit(10)
		}

		printReqErr := output.PrintRequest(req)
		if printReqErr != nil {
			fmt.Println("Error while printing request.")
			fmt.Println(printReqErr)
			os.Exit(20)
		}

		client := &http.Client{}
		resp, respErr := client.Do(req)

		if respErr != nil {
			fmt.Println("There was an error.")
			fmt.Println(respErr)
			os.Exit(30)
		}

		printRespErr := output.PrintResponse(resp)
		if printRespErr != nil {
			fmt.Println("Error while printing response.")
			fmt.Println(printRespErr)
			os.Exit(40)
		}
	*/
}
