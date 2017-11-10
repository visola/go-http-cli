package strings

import (
	"strings"
)

// ParseExpression parses a string template with variables in the format ${myVariable} and uses
// values from the context map to replace them. The resulting string is the template string with
// replaced values.
func ParseExpression(template string, context map[string]string) string {
	result := template
	for key, value := range context {
		result = strings.Replace(result, "{"+key+"}", value, -1)
	}
	return result
}
