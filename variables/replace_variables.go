package variables

import (
	"fmt"
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
			fmt.Printf("Last position is: %d\n", lastPosition)
			fmt.Printf("Variable starts at: %d, ends at: %d\n", variable.Start, variable.End)
			if lastPosition != variable.Start {
				fmt.Println(template[lastPosition:variable.Start])
				resultBuilder.WriteString(template[lastPosition:variable.Start])
			}
			resultBuilder.WriteString(value)
			lastPosition = variable.End
			fmt.Printf("String so far: '%s'\n", resultBuilder.String())
		}
	}

	resultBuilder.WriteString(template[lastPosition:])
	return resultBuilder.String()
}
