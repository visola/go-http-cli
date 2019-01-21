package request

import "strings"

func getContentType(headers map[string][]string) string {
	for name, values := range headers {
		if strings.ToLower(strings.TrimSpace(name)) == "content-type" {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}
