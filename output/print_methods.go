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
	if response.StatusCode < 300 {
		color.Green("%s %s\n", response.Status, response.Protocol)
	} else if response.StatusCode < 400 {
		color.Yellow("%s %s\n", response.Status, response.Protocol)
	} else {
		color.Red("%s %s\n", response.Status, response.Protocol)
	}

	printHeaders(response.Headers)
	printBody(response.Body, "<<")
}

func printBody(body string, linePrefix string) {
	if body != "" {
		bodyColor := color.New(color.Bold).PrintfFunc()
		split := strings.Split(body, "\n")
		for _, line := range split {
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
	sentHeaderKeyColor := color.New(color.Bold, color.FgBlack).PrintfFunc()
	sentHeaderValueColor := color.New(color.FgBlack).PrintfFunc()

	for headerName, values := range headers {
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
}
