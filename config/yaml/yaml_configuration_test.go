package yaml

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBaseURL     = "http://www.test.com"
	testHeaderName  = "Content-type"
	testHeaderValue = "application/json"
)

func TestReadFrom(t *testing.T) {
	t.Run("Parses BaseURL correctly", testReadBaseURL)
	t.Run("Parses Headers when passed as Array", testReadHeadersAsArray)
	t.Run("Parses Headers when passed as String", testReadHeadersAsString)
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

	assertNoError(t, err)
	assert.Equal(t, len(parsedConfiguration.Headers()), 1, "Should have parsed one header")
	assert.Equal(t, parsedConfiguration.Headers()[testHeaderName], []string{testHeaderValue}, "Should parse header value correctly")
}

func testReadHeadersAsString(t *testing.T) {
	yamlContent := "headers:" + "\n  " + testHeaderName + ": " + testHeaderValue
	file := writeToFile("test.yml", yamlContent)
	parsedConfiguration, err := ReadFrom(file.Name())

	assertNoError(t, err)
	assert.Equal(t, len(parsedConfiguration.Headers()), 1, "Should have parsed one header")
	assert.Equal(t, parsedConfiguration.Headers()[testHeaderName], []string{testHeaderValue}, "Should parse header value correctly")
}

func assertNoError(t *testing.T, err error) {
	assert.Nil(t, err, "Should not return error")
	if err != nil {
		return
	}
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
