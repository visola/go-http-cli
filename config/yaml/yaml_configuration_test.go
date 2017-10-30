package yaml

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBaseURL       = "http://www.test.com"
	testHeaderName    = "Content-type"
	testHeaderValue   = "application/json"
	testVariableName  = "someVariable"
	testVariableValue = "someValue"
)

func TestReadFrom(t *testing.T) {
	t.Run("Parses BaseURL correctly", testReadBaseURL)
	t.Run("Parses Headers when passed as Array", testReadHeadersAsArray)
	t.Run("Parses Headers when passed as String", testReadHeadersAsString)
	t.Run("Parses Variables when passed as Array", testReadVariablesAsArray)
	t.Run("Parses Variables when passed as String", testReadVariablesAsString)
}

func testReadBaseURL(t *testing.T) {
	yamlContent := "baseURL: " + testBaseURL
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())

	assertNoError(t, err)
	assert.Equal(t, testBaseURL, parsedConfiguration.BaseURL())
}

func testReadHeadersAsArray(t *testing.T) {
	yamlContent := "headers:" + "\n  " + testHeaderName + ":\n  - " + testHeaderValue
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())
	assertParsedHeaders(t, err, parsedConfiguration)
}

func testReadHeadersAsString(t *testing.T) {
	yamlContent := "headers:" + "\n  " + testHeaderName + ": " + testHeaderValue
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())
	assertParsedHeaders(t, err, parsedConfiguration)
}

func testReadVariablesAsArray(t *testing.T) {
	yamlContent := "variables:" + "\n  " + testVariableName + ":\n  - " + testVariableValue
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())
	assertParsedVariables(t, err, parsedConfiguration)
}

func testReadVariablesAsString(t *testing.T) {
	yamlContent := "variables:" + "\n  " + testVariableName + ": " + testVariableValue
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())
	assertParsedVariables(t, err, parsedConfiguration)
}

func assertNoError(t *testing.T, err error) {
	assert.Nil(t, err, "Should not return error")
	if err != nil {
		return
	}
}

func assertParsedHeaders(t *testing.T, err error, parsedConfiguration *FileConfiguration) {
	assertNoError(t, err)
	assert.Equal(t, len(parsedConfiguration.Headers()), 1, "Should have parsed one header")
	assert.Equal(t, parsedConfiguration.Headers()[testHeaderName], []string{testHeaderValue}, "Should parse header value correctly")
}

func assertParsedVariables(t *testing.T, err error, parsedConfiguration *FileConfiguration) {
	assertNoError(t, err)
	assert.Equal(t, len(parsedConfiguration.Variables()), 1, "Should have parsed one header")
	assert.Equal(t, parsedConfiguration.Variables()[testVariableName], []string{testVariableValue}, "Should parse variable value correctly")
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
