package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestParse(t *testing.T) {
	t.Run("Parses profile correctly", testParsesProfileCorrectly)
	t.Run("Fails if profile file does not exist", testFailIfProfileFileDoesNotExist)
	t.Run("Failes if configuration file does not exist", testConfigurationFileDoesNotExist)
	t.Run("Handles failure to parse Yaml file", testFailToParseYamlFile)
}

func testParsesProfileCorrectly(t *testing.T) {
	profileName := "profile"
	tmpFile := writeToFile(profileName+".yaml", "headers:\n "+testHeader+": "+testValue)
	os.Setenv("GO_HTTP_PROFILES", filepath.Dir(tmpFile.Name()))
	args := []string{"--method", testMethod, "--data", testData, "+" + profileName, testURL}
	configuration, err := Parse(args)
	assertCorrectlyParsed(t, configuration, err)
	os.Unsetenv("GO_HTTP_PROFILES")
}

func testFailIfProfileFileDoesNotExist(t *testing.T) {
	args := []string{"+profile", testURL}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func testConfigurationFileDoesNotExist(t *testing.T) {
	args := []string{"--method", testMethod, "--data", testData, "--config", "fileThatDoesNotExist.yml", testURL}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func testFailToParseYamlFile(t *testing.T) {
	tmpFile := writeToFile("some_file.yml", "bla bla")
	args := []string{"--method", testMethod, "--data", testData, "--config", tmpFile.Name(), testURL}
	_, err := Parse(args)
	assert.NotNil(t, err, "Should return error")
}

func assertCorrectlyParsed(t *testing.T, configuration Configuration, err error) {
	assert.Nil(t, err, "Should not return error")
	if configuration == nil {
		return
	}
	assert.Equal(t, 1, len(configuration.Headers()), "Should parse one header correctly")
	assert.Equal(t, []string{testValue}, configuration.Headers()[testHeader], "Should parse the correct value for the header")
	assert.Equal(t, testMethod, configuration.Method(), "Should parse method correctly")
	assert.Equal(t, testURL, configuration.URL(), "Should parse URL correctly")
	assert.Equal(t, testData, configuration.Body(), "Should parse data correctly")
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
