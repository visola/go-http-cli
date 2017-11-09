package output

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

// PrintResponse outputs a http.Response
func PrintResponse(response *http.Response) error {
	defer response.Body.Close()

	color.Green("\n%s\n", response.Status)

	receivedHeaderKeyColor := color.New(color.Bold, color.FgBlack).PrintfFunc()
	receivedHeaderValueColor := color.New(color.FgBlack).PrintfFunc()
	for k, vs := range response.Header {
		receivedHeaderKeyColor("%s:", k)
		receivedHeaderValueColor(" %s\n", strings.Join(vs, ", "))
	}

	bodyBytes, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return readErr
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

	return nil
}
