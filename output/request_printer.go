package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/request"
)

type bodyBuffer struct {
	*strings.Reader
}

func (bb *bodyBuffer) Close() error {
	return nil
}

// PrintRequest outputs the http.Request
func PrintRequest(request request.Request) error {
	color.Green("\n%s %s\n", request.Method, request.URL)

	sentHeaderKeyColor := color.New(color.Bold, color.FgBlue).PrintfFunc()
	sentHeaderValueColor := color.New(color.FgBlue).PrintfFunc()

	for headerName, values := range request.Headers {
		sentHeaderKeyColor("%s:", headerName)
		if len(values) > 1 {
			for _, val := range values {
				sentHeaderValueColor("\n  %s", val)
			}
			fmt.Println("")
		} else {
			sentHeaderValueColor(" %s\n", values[0])
		}
	}

	if request.Body != "" {
		split := strings.Split(request.Body, "\n")
		for _, line := range split {
			fmt.Printf(">> %s\n", line)
		}
	}

	return nil
}
