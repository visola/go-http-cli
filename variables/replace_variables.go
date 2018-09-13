package variables

import (
	"strings"
)

// ReplaceVariables replace all variables found in template matching the names
// from context
func ReplaceVariables(template string, context map[string]string) string {
	variables := FindVariables(template)

	var resultBuilder strings.Builder
	lastPosition := 0
	for _, variable := range variables {
		value, ok := context[variable.Name]
		if ok { // variable in context
			if lastPosition != variable.Start {
				resultBuilder.WriteString(template[lastPosition:variable.Start])
			}
			resultBuilder.WriteString(value)
			lastPosition = variable.End
		}
	}

	resultBuilder.WriteString(template[lastPosition:])
	return resultBuilder.String()
}
