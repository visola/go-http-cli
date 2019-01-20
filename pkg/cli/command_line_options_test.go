package cli

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBaseURL = "http://www.test.com"
	testData    = "Some Data"
	testMethod  = http.MethodPost
	testHeader  = "SomeHeader"
	testValue   = "SomeValue"
	testURL     = "/somePath"
)

func TestParseCommandLineOptions(t *testing.T) {
	t.Run("Parses a full URL correctly", testParsesFullURLCorrectly)
	t.Run("Parses a full URL with query string correctly", testParsesFullURLWithQueryStringCorrectly)

	t.Run("Parses a path correctly", testParsesPath)
	t.Run("Parses a path with query string", testParsesPathWithQueryString)

	t.Run("Parses values correctly", testParsesValuesCorrectly)

	t.Run("Parses all arguments using short names", testParsesShortNames)
	t.Run("Parses all arguments using long names", testParsesLongNames)

	t.Run("Parses multiple values for the same header", testParsesMultipleValuesForHeader)
	t.Run("Parses header with = on the value", testParsesHeaderWithEqualOnValue)
	t.Run("Fails to parse header with wrong separator", testFailToParseHeaderWithWrongSeparator)
}

func testParsesFullURLCorrectly(t *testing.T) {
	url := testBaseURL + testURL
	args := []string{url}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, url, configuration.URL, "Should parse URL correctly")
}

func testParsesFullURLWithQueryStringCorrectly(t *testing.T) {
	url := testBaseURL + testURL + "?someKey=someValue"
	args := []string{url}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, url, configuration.URL, "Should parse URL correctly")
}

func testParsesPath(t *testing.T) {
	url := testURL
	args := []string{url}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, url, configuration.URL, "Should parse URL correctly")
}

func testParsesPathWithQueryString(t *testing.T) {
	url := testURL + "?someKey=someValue"
	args := []string{url}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, url, configuration.URL, "Should parse URL correctly")
}

func testParsesValuesCorrectly(t *testing.T) {
	value := "value with spaces"
	key := "key$%bla"
	args := []string{key + "=" + value}
	configuration, err := ParseCommandLineOptions(args)

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, 1, len(configuration.Values), "Parsed one key")
	assert.Equal(t, 1, len(configuration.Values[key]), "Parsed one value")
	assert.Equal(t, value, configuration.Values[key][0], "Parses values correctly")
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
