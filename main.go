package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/visola/go-http-cli/config"
)

type bodyBuffer struct {
	*bytes.Buffer
}

func (bb *bodyBuffer) Close() error {
	return nil
}

func main() {
	configuration, err := config.Parse(os.Args[1:])

	if err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}

	url := configuration.URL()
	if baseURL := configuration.BaseURL(); baseURL != "" {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL = baseURL + "/"
		}

		if strings.HasPrefix(url, "/") {
			url = url[1:]
		}

		url = baseURL + url
	}

	color.Green("\n%s %s\n", configuration.Method(), url)
	req, reqErr := http.NewRequest(configuration.Method(), url, nil)

	if reqErr != nil {
		color.Red("Error while creating request: %s", reqErr)
		os.Exit(10)
	}

	for k, vs := range configuration.Headers() {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	sentHeaderKeyColor := color.New(color.Bold, color.FgBlue).PrintfFunc()
	sentHeaderValueColor := color.New(color.FgBlue).PrintfFunc()
	for k, vs := range req.Header {
		sentHeaderKeyColor("%s:", k)
		sentHeaderValueColor(" %s\n", strings.Join(vs, ", "))
	}

	if configuration.Body() != "" {
		split := strings.Split(configuration.Body(), "\n")
		for _, line := range split {
			fmt.Printf(">> %s\n", line)
		}

		req.Body = &bodyBuffer{bytes.NewBufferString(configuration.Body())}
	}

	client := &http.Client{}
	resp, respErr := client.Do(req)

	if respErr != nil {
		fmt.Println("There was an error.")
		fmt.Println(respErr)
		os.Exit(20)
	}

	defer resp.Body.Close()

	color.Green("\n%s\n", resp.Status)

	receivedHeaderKeyColor := color.New(color.Bold, color.FgBlack).PrintfFunc()
	receivedHeaderValueColor := color.New(color.FgBlack).PrintfFunc()
	for k, vs := range resp.Header {
		receivedHeaderKeyColor("%s:", k)
		receivedHeaderValueColor(" %s\n", strings.Join(vs, ", "))
	}

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		color.Red("Error while reading body. %s", readErr)
		os.Exit(30)
	}

	if len(bodyBytes) != 0 {
		split := strings.Split(string(bodyBytes), "\n")
		fmt.Println("")

		receivedBodyColor := color.New(color.Bold).PrintfFunc()
		for _, line := range split {
			receivedBodyColor("%s\n", line)
		}
		fmt.Println("")
	}

}
