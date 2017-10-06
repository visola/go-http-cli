package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	data   = "Some Data"
	method = "POST"
	header = "SomeHeader"
	value  = "SomeValue"
	url    = "http://www.google.com"
)

func TestParse(t *testing.T) {
	t.Run("Parses all arguments using short names", testParsesShortNames)
	t.Run("Parses all arguments using long names", testParsesLongNames)
	t.Run("Parses multiple values for the same header", testParsesMultipleValuesForHeader)
	t.Run("Parses configuration from file", testParsesConfigurationFromFile)
	t.Run("Fails to parse header with wrong separator", testFailToParseHeaderWithWrongSeparator)
}

func testParsesShortNames(t *testing.T) {
	args := []string{"--method", method, "--data", data, "--header", header + "=" + value, url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesLongNames(t *testing.T) {
	args := []string{"-X", method, "-d", data, "-H", header + "=" + value, url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesMultipleValuesForHeader(t *testing.T) {
	newValue := "AnotherValue"
	args := []string{"--header", header + "=" + value, "--header", header + "=" + newValue, url}
	configuration, err := Parse(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{value, newValue}, configuration.Headers()[header], "Should parse correct values for header")
	assert.Equal(t, url, configuration.URL(), "Should parse URL correctly")
}

func testFailToParseHeaderWithWrongSeparator(t *testing.T) {
	args := []string{"--header", header + ":" + value, url}
	_, err := Parse(args)

	assert.NotNil(t, err, "Should return error")
	assert.Regexp(t, "^Error while parsing header", err.Error())
}

func testParsesConfigurationFromFile(t *testing.T) {
	simpleHeaderYaml := "headers:\n  " + header + ":\n    - " + value
	tmpFile, err := ioutil.TempFile("", "simple_header.yml")

	if err != nil {
		panic(err)
	}

	tmpFile.WriteString(simpleHeaderYaml)

	args := []string{"--method", method, "--data", data, "--config", tmpFile.Name(), url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
}

func assertCorrectlyParsed(t *testing.T, configuration Configuration, err error) {
	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{value}, configuration.Headers()[header], "Should parse the correct value for the header")
	assert.Equal(t, method, configuration.Method(), "Should parse method correctly")
	assert.Equal(t, url, configuration.URL(), "Should parse URL correctly")
	assert.Equal(t, data, configuration.Body(), "Should parse data correctly")
}
