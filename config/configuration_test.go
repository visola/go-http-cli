package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	t.Run("Parses profile correctly", testParsesProfileCorrectly)
	t.Run("Fails if profile file does not exist", testFailIfProfileFileDoesNotExist)
	t.Run("Parses configuration from file using header as string", testParsesConfigurationFromFileUsingHeaderString)
	t.Run("Parses configuration from file using header as array", testParsesConfigurationFromFileUsingHeaderAsArray)
	t.Run("Fails to parse header with wrong separator", testFailToParseHeaderWithWrongSeparator)
	t.Run("Failes if configuration file does not exist", testConfigurationFileDoesNotExist)
	t.Run("Handles failure to parse Yaml file", testFailToParseYamlFile)
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

func testParsesProfileCorrectly(t *testing.T) {
	profileName := "profile"
	tmpFile := writeToFile(profileName+".yaml", "headers:\n "+header+": "+value)
	os.Setenv("GO_HTTP_PROFILES", filepath.Dir(tmpFile.Name()))
	args := []string{"--method", method, "--data", data, "+" + profileName, url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
	os.Unsetenv("GO_HTTP_PROFILES")
}

func testFailIfProfileFileDoesNotExist(t *testing.T) {
	args := []string{"+profile", url}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func testParsesConfigurationFromFileUsingHeaderString(t *testing.T) {
	tmpFile := writeToFile("some_file.yml", "headers:\n  "+header+": "+value)
	args := []string{"--method", method, "--data", data, "--config", tmpFile.Name(), url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testParsesConfigurationFromFileUsingHeaderAsArray(t *testing.T) {
	tmpFile := writeToFile("some_file.yml", "headers:\n  "+header+":\n    - "+value)
	args := []string{"--method", method, "--data", data, "--config", tmpFile.Name(), url}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
}

func testConfigurationFileDoesNotExist(t *testing.T) {
	args := []string{"--method", method, "--data", data, "--config", "fileThatDoesNotExist.yml", url}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func testFailToParseYamlFile(t *testing.T) {
	tmpFile := writeToFile("some_file.yml", "bla bla")
	args := []string{"--method", method, "--data", data, "--config", tmpFile.Name(), url}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func assertCorrectlyParsed(t *testing.T, configuration Configuration, err error) {
	assert.Nil(t, err, "Should not return error")
	if configuration == nil {
		return
	}
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{value}, configuration.Headers()[header], "Should parse the correct value for the header")
	assert.Equal(t, method, configuration.Method(), "Should parse method correctly")
	assert.Equal(t, url, configuration.URL(), "Should parse URL correctly")
	assert.Equal(t, data, configuration.Body(), "Should parse data correctly")
}

func writeToFile(fileName string, content string) *os.File {
	tempDir, createDirErr := ioutil.TempDir("", "test")

	if createDirErr != nil {
		panic(createDirErr)
	}

	tmpFile, err := os.Create(tempDir + "/" + fileName)
	if err != nil {
		panic(err)
	}

	_, writeErr := tmpFile.WriteString(content)
	if writeErr != nil {
		panic(writeErr)
	}

	return tmpFile
}
