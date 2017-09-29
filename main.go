package main

import (
	"bytes"
	"fmt"
	"github.com/visola/go-http-cli/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type bodyBuffer struct {
	*bytes.Buffer
}

func (bb *bodyBuffer) Close() error {
	return nil
}

func main() {
	configuration, err := config.Parse()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if configuration.Url == "" {
		fmt.Println("Nothing to do.")
		os.Exit(2)
	}

	fmt.Printf("\n%s %s\n", configuration.Method, configuration.Url)

	if len(configuration.Headers) == 0 {
		fmt.Println(">>")
	} else {
		for k, v := range configuration.Headers {
			fmt.Printf(">> '%s' = '%s'\n", k, v)
		}
	}

	req, reqErr := http.NewRequest(configuration.Method, configuration.Url, nil)

	if reqErr != nil {
		fmt.Println("Error while creating request: ", reqErr)
		os.Exit(10)
	}

	for k, v := range configuration.Headers {
		req.Header.Add(k, v)
	}

	if configuration.Body != "" {
		fmt.Println(">>")
		split := strings.Split(configuration.Body, "\n")
		for _, line := range split {
			fmt.Printf(">> %s\n", line)
		}

		req.Body = &bodyBuffer{bytes.NewBufferString(configuration.Body)}
	}

	fmt.Println("--")

	client := &http.Client{}
	resp, respErr := client.Do(req)

	if respErr != nil {
		fmt.Println("There was an error.")
		fmt.Println(respErr)
		os.Exit(20)
	}

	defer resp.Body.Close()

	for k, v := range resp.Header {
		fmt.Printf("<< '%s' = '%s'\n", k, v)
	}

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error while reading body.", readErr)
		os.Exit(30)
	}

	fmt.Println("<<")
	if len(bodyBytes) != 0 {
		split := strings.Split(string(bodyBytes), "\n")
		for _, line := range split {
			fmt.Printf("<< %s\n", line)
		}
		fmt.Println("<<")
	}

}
