package request

import (
	"strings"

	myStrings "github.com/visola/go-http-cli/strings"
)

// ParseURL parses the configuration to generate the final URL that the request will be sent to.
func ParseURL(base string, path string, variables map[string]string) string {
	url := path

	if !strings.HasPrefix(url, "http") && base != "" {
		if !strings.HasSuffix(base, "/") {
			base = base + "/"
		}

		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}

		url = base + path
	}

	return myStrings.ParseExpression(url, variables)
}
