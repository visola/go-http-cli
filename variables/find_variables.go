package variables

import (
	"regexp"
)

var variableRegexp, _ = regexp.Compile("{(\\w+){1}(?::(\\w+)){0,1}}")

// FindVariables find variables with tag (if any) in a given string
func FindVariables(stringWithVariables string) []Variable {
	result := make([]Variable, 0)

	matches := variableRegexp.FindAllStringSubmatch(stringWithVariables, -1)
	allIndexes := variableRegexp.FindAllStringSubmatchIndex(stringWithVariables, -1)
	for matchCount, match := range matches {
		matchIndexes := allIndexes[matchCount]
		result = append(result, Variable{
			End:       matchIndexes[1],
			Name:      match[1],
			NameEnd:   matchIndexes[3],
			NameStart: matchIndexes[2],
			Start:     matchIndexes[0],
			Tag:       match[2],
			TagEnd:    matchIndexes[5],
			TagStart:  matchIndexes[4],
		})
	}

	return result
}
