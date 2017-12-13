package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBaseURL = "http://www.test.com"
	testData    = "Some Data"
	testMethod  = "POST"
	testHeader  = "SomeHeader"
	testValue   = "SomeValue"
	testURL     = "/somePath"
)

func TestParseCommandLineOptions(t *testing.T) {
	t.Run("Uses POST method if no method given and data present", testAddsDefaultMethodWithData)
	t.Run("Uses GET method if no method given and no data present", testAddsDefaultMethodWithoutData)
	t.Run("Parses all arguments using short names", testParsesShortNames)
	t.Run("Parses all arguments using long names", testParsesLongNames)
	t.Run("Parses multiple values for the same header", testParsesMultipleValuesForHeader)
	t.Run("Parses header with = on the value", testParsesHeaderWithEqualOnValue)
	t.Run("Fails to parse header with wrong separator", testFailToParseHeaderWithWrongSeparator)
}

func testAddsDefaultMethodWithoutData(t *testing.T) {
	args := []string{"--header", testHeader + "=" + testValue, testURL}
	configuration, err := ParseCommandLineOptions(args)
	assert.Nil(t, err, "Should not return error")
	if configuration == nil {
		return
	}
	assert.Equal(t, 1, len(configuration.Headers), "Should parse one header correctly")
	assert.Equal(t, []string{testValue}, configuration.Headers[testHeader], "Should parse the correct value for the header")
	assert.Equal(t, "GET", configuration.Method, "Should parse method correctly")
	assert.Equal(t, testURL, configuration.URL, "Should parse URL correctly")
	assert.Equal(t, "", configuration.Body, "Should parse data correctly")
}

func testAddsDefaultMethodWithData(t *testing.T) {
	args := []string{"--data", testData, "--header", testHeader + "=" + testValue, testURL}
	configuration, err := ParseCommandLineOptions(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesShortNames(t *testing.T) {
	args := []string{"--method", testMethod, "--data", testData, "--header", testHeader + "=" + testValue, testURL}
	configuration, err := ParseCommandLineOptions(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesLongNames(t *testing.T) {
	args := []string{"-X", testMethod, "-d", testData, "-H", testHeader + "=" + testValue, testURL}
	configuration, err := ParseCommandLineOptions(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesMultipleValuesForHeader(t *testing.T) {
	newValue := "AnotherValue"
	args := []string{"--header", testHeader + ":" + testValue, "--header", testHeader + "=" + newValue, testURL}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers), "Should parse one header correctly")
	assert.Equal(t, []string{testValue, newValue}, configuration.Headers[testHeader], "Should parse correct values for header")
	assert.Equal(t, testURL, configuration.URL, "Should parse URL correctly")
}

func testParsesHeaderWithEqualOnValue(t *testing.T) {
	newValue := "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8"
	args := []string{"--header", testHeader + "=" + newValue, testURL}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Headers), "Should parse one header correctly")
	assert.Equal(t, []string{newValue}, configuration.Headers[testHeader], "Should parse correct value for header")
	assert.Equal(t, testURL, configuration.URL, "Should parse URL correctly")
}

func testFailToParseHeaderWithWrongSeparator(t *testing.T) {
	args := []string{"--header", testHeader + "->" + testValue, testURL}
	_, err := ParseCommandLineOptions(args)

	assert.NotNil(t, err, "Should return error")
	if t != nil {
		assert.Regexp(t, "^Error while parsing key value pair", err.Error())
	}
}

func assertCorrectlyParsed(t *testing.T, configuration *CommandLineOptions, err error) {
	assert.Nil(t, err, "Should not return error")
	assert.NotNil(t, configuration, "Should return a configuration")
	assert.Equal(t, 1, len(configuration.Headers), "Should parse one header correctly")
	assert.Equal(t, []string{testValue}, configuration.Headers[testHeader], "Should parse the correct value for the header")
	assert.Equal(t, testMethod, configuration.Method, "Should parse method correctly")
	assert.Equal(t, testURL, configuration.URL, "Should parse URL correctly")
	assert.Equal(t, testData, configuration.Body, "Should parse data correctly")
}
