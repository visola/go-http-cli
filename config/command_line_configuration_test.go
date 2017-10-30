package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommandLine(t *testing.T) {
	t.Run("Uses POST method if no method given and data present", testAddsDefaultMethodWithData)
	t.Run("Uses GET method if no method given and no data present", testAddsDefaultMethodWithoutData)
	t.Run("Parses all arguments using short names", testParsesShortNames)
	t.Run("Parses all arguments using long names", testParsesLongNames)
	t.Run("Parses multiple values for the same header", testParsesMultipleValuesForHeader)
	t.Run("Parses profile correctly", testParsesProfileCorrectly)
	t.Run("Fails to parse header with wrong separator", testFailToParseHeaderWithWrongSeparator)
}

func testAddsDefaultMethodWithoutData(t *testing.T) {
	args := []string{"--header", testHeader + "=" + testValue, testURL}
	configuration, err := parseCommandLine(args)
	assert.Nil(t, err, "Should not return error")
	if configuration == nil {
		return
	}
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{testValue}, configuration.Headers()[testHeader], "Should parse the correct value for the header")
	assert.Equal(t, "GET", configuration.Method(), "Should parse method correctly")
	assert.Equal(t, testURL, configuration.URL(), "Should parse URL correctly")
	assert.Equal(t, "", configuration.Body(), "Should parse data correctly")
}

func testAddsDefaultMethodWithData(t *testing.T) {
	args := []string{"--data", testData, "--header", testHeader + "=" + testValue, testURL}
	configuration, err := parseCommandLine(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesShortNames(t *testing.T) {
	args := []string{"--method", testMethod, "--data", testData, "--header", testHeader + "=" + testValue, testURL}
	configuration, err := parseCommandLine(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesLongNames(t *testing.T) {
	args := []string{"-X", testMethod, "-d", testData, "-H", testHeader + "=" + testValue, testURL}
	configuration, err := parseCommandLine(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesMultipleValuesForHeader(t *testing.T) {
	newValue := "AnotherValue"
	args := []string{"--header", testHeader + "=" + testValue, "--header", testHeader + "=" + newValue, testURL}
	configuration, err := parseCommandLine(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{testValue, newValue}, configuration.Headers()[testHeader], "Should parse correct values for header")
	assert.Equal(t, testURL, configuration.URL(), "Should parse URL correctly")
}

func testFailToParseHeaderWithWrongSeparator(t *testing.T) {
	args := []string{"--header", testHeader + ":" + testValue, testURL}
	_, err := parseCommandLine(args)

	assert.NotNil(t, err, "Should return error")
	assert.Regexp(t, "^Error while parsing key value pair", err.Error())
}
