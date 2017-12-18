package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/request"
)

// PrintResponse outputs a http.Response
func PrintResponse(response *request.HTTPResponse) error {
	color.Green("\n%s %s\n", response.Status, response.Protocol)

	receivedHeaderKeyColor := color.New(color.Bold, color.FgBlack).PrintfFunc()
	receivedHeaderValueColor := color.New(color.FgBlack).PrintfFunc()

	for headerName, values := range response.Headers {
		receivedHeaderKeyColor("%s:", headerName)
		if len(values) > 1 {
			for _, val := range values {
				receivedHeaderValueColor("\n  %s", val)
			}
			fmt.Println("")
		} else {
			receivedHeaderValueColor(" %s\n", values[0])
		}
	}

	if len(response.Body) != 0 {
		split := strings.Split(response.Body, "\n")
		fmt.Println("")

		receivedBodyColor := color.New(color.Bold).PrintfFunc()
		for _, line := range split {
			receivedBodyColor("%s\n", line)
		}
		fmt.Println("")
	}

	return nil
}
