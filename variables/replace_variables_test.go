package variables

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceVariables(t *testing.T) {
	template := `
{name} was born in {dob} in {location}.
{name} hair is {hairColor}.
{this} is a variable that is not in the context.
Two variables togethers should also work: {dob}{location}.
`

	context := map[string]string{
		"dob":       "01/08/1956",
		"hairColor": "black",
		"location":  "Japan",
		"name":      "John",
	}

	assert.Equal(
		t,
		`
John was born in 01/08/1956 in Japan.
John hair is black.
{this} is a variable that is not in the context.
Two variables togethers should also work: 01/08/1956Japan.
`,
		ReplaceVariables(template, context),
		"Should find and replace available variables correctly",
	)
}
