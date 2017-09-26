package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type headerFlags []string

func (i *headerFlags) String() string {
	return "No String Representation"
}

func (i *headerFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var method string
	var headers headerFlags

	flag.StringVar(&method, "method", "GET", "HTTP method to be used")
	flag.Var(&headers, "header", "Headers to include with your request")

	flag.Parse()

	fmt.Println("Method: ", method)

	if len(flag.Args()) != 1 {
		fmt.Println("Nothing to do.")
		os.Exit(1)
	}

	url := flag.Args()[0]
	fmt.Println("URL: ", url)

	req, reqErr := http.NewRequest("GET", url, nil)

	if reqErr != nil {
		fmt.Println("Error while creating request: ", reqErr)
		os.Exit(10)
	}

	for _, kv := range headers {
		s := strings.Split(kv, "=")
		if len(s) != 2 {
			fmt.Println("Error while parsing header: ", kv)
			fmt.Println("Should be a '=' separated key/value, e.g.: Content-type=application/x-www-form-urlencoded")
			os.Exit(11)
		}
		req.Header.Add(s[0], s[1])
	}

	client := &http.Client{}
	resp, respErr := client.Do(req)

	if respErr != nil {
		fmt.Println("There was an error.")
		fmt.Println(respErr)
		os.Exit(20)
	}

	defer resp.Body.Close()

	for k, v := range req.Header {
		fmt.Printf(">> '%s' = '%s'\n", k, v)
	}
	fmt.Println("---")

	for k, v := range resp.Header {
		fmt.Printf("<< '%s' = '%s'\n", k, v)
	}

	fmt.Println("---")

	bodyBytes, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		fmt.Println("Error while reading body.", readErr)
		os.Exit(30)
	}

	fmt.Println(string(bodyBytes))

}
