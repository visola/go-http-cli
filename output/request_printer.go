package output

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

type bodyBuffer struct {
	*strings.Reader
}

func (bb *bodyBuffer) Close() error {
	return nil
}

// PrintRequest outputs the http.Request
func PrintRequest(request *http.Request) error {
	color.Green("\n%s %s\n", request.Method, request.URL)

	sentHeaderKeyColor := color.New(color.Bold, color.FgBlue).PrintfFunc()
	sentHeaderValueColor := color.New(color.FgBlue).PrintfFunc()
	for k, vs := range request.Header {
		sentHeaderKeyColor("%s:", k)
		sentHeaderValueColor(" %s\n", strings.Join(vs, ", "))
	}

	if request.Body != nil {
		defer request.Body.Close()
		bytes, readError := ioutil.ReadAll(request.Body)
		if readError != nil {
			return readError
		}

		body := string(bytes)

		// rewind reader
		request.Body = &bodyBuffer{strings.NewReader(body)}

		split := strings.Split(body, "\n")
		for _, line := range split {
			fmt.Printf(">> %s\n", line)
		}
	}

	return nil
}
