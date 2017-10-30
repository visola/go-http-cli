package strings

import (
	"fmt"
)

func ExampleParseExpression() {
	template := "${name} was born in ${dob} in ${location}. ${name} hair is ${hairColor}."
	context := map[string]string{
		"dob":       "01/08/1956",
		"hairColor": "black",
		"location":  "Japan",
		"name":      "John",
	}
	fmt.Println(ParseExpression(template, context))
	// Output: John was born in 01/08/1956 in Japan. John hair is black.
}
