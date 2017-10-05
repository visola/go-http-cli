package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var args []string
	var configuration Configuration
	var err error

	data := "Some Data"
	method := "POST"
	header := "SomeHeader"
	value := "SomeValue"
	url := "http://www.google.com"

	// Parse all arguments correctly using long format
	args = []string{"--method", method, "--data", data, "--header", header + "=" + value, url}
	configuration, err = Parse(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{value}, configuration.Headers()[header], "Should parse the correct value for the header")
	assert.Equal(t, data, configuration.Body(), "Should parse data correctly")

	// Parse all arguments correctly using short format
	args = []string{"-X", method, "-d", data, "-H", header + "=" + value, url}
	configuration, err = Parse(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{value}, configuration.Headers()[header], "Should parse the correct value for the header")
	assert.Equal(t, data, configuration.Body(), "Should parse data correctly")

	// Fail to parse header with wrong separator
	args = []string{"--header", header + ":" + value, url}
	configuration, err = Parse(args)

	assert.NotNil(t, err, "Should return error")
	assert.Regexp(t, "^Error while parsing header", err.Error())
}
