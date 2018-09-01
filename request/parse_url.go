package request

import (
	"strings"
)

// ParseURL parses the configuration to generate the final URL that the request will be sent to
func ParseURL(base string, paths ...string) string {
	path := coalesce(paths...)
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

	return url
}

func coalesce(values ...string) string {
	for _, oneValue := range values {
		if oneValue != "" {
			return oneValue
		}
	}

	return ""
}
