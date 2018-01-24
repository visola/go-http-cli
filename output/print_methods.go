package output

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/request"
)

// PrintRequest outputs the http.Request
func PrintRequest(request request.Request) {
	boldGreen := color.New(color.Bold, color.FgGreen).PrintfFunc()
	boldGreen("\n%s %s\n", request.Method, request.URL)

	printHeaders(request.Headers)
	printCookies(request.Cookies)
	printBody(request.Body, ">>")
}

// PrintResponse outputs a http.Response
func PrintResponse(response request.Response) {
	printSummaryFunction := color.Green
	if response.StatusCode >= 300 && response.StatusCode < 400 {
		printSummaryFunction = color.Yellow
	} else {
		printSummaryFunction = color.Red
	}

	printSummaryFunction("%s %s\n", response.Status, response.Protocol)
	printHeaders(response.Headers)
	printBody(response.Body, "<<")
}

func printBody(body string, linePrefix string) {
	if body != "" {
		bodyColor := color.New(color.Bold).PrintfFunc()
		for _, line := range strings.Split(body, "\n") {
			bodyColor("%s %s\n", linePrefix, line)
		}
	}
}

func printCookies(cookies []*http.Cookie) {
	if len(cookies) > 0 {
		sentCookieKeyColor := color.New(color.Bold, color.FgBlue).PrintfFunc()
		sentCookieValueColor := color.New(color.FgBlue).PrintfFunc()

		sentCookieKeyColor("Cookies:")
		for _, cookie := range cookies {
			sentCookieKeyColor("\n  %s: ", cookie.Name)
			sentCookieValueColor("%s", cookie.Value)
		}
		fmt.Println("")
	}
}

func printHeaders(headers map[string][]string) {
	headerKeyColor := color.New(color.Bold, color.FgBlack).PrintfFunc()
	headerValueColor := color.New(color.FgBlack).PrintfFunc()

	for headerName, values := range headers {
		headerKeyColor("%s:", headerName)
		if len(values) == 1 {
			headerValueColor(" %s\n", values[0])
		} else if len(values) > 1 {
			for _, val := range values {
				headerValueColor("\n  %s", val)
			}
			fmt.Println("")
		}
	}
}
