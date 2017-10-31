package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfiguration struct {
	baseURL   string
	headers   map[string][]string
	body      string
	method    string
	url       string
	variables map[string]string
}

func (conf testConfiguration) BaseURL() string {
	return conf.baseURL
}

func (conf testConfiguration) Headers() map[string][]string {
	return conf.headers
}

func (conf testConfiguration) Body() string {
	return conf.body
}

func (conf testConfiguration) Method() string {
	return conf.method
}

func (conf testConfiguration) URL() string {
	return conf.url
}

func (conf testConfiguration) Variables() map[string]string {
	return conf.variables
}

func TestHierarchicalConfiguration(t *testing.T) {
	t.Run("Overrides base url correctly", testOverrideBaseURL)
	t.Run("Overrides headers correctly", testOverrideHeaders)
	t.Run("Overrides variables correctly", testOverrideVariables)
}

func testOverrideBaseURL(t *testing.T) {
	config1 := testConfiguration{
		baseURL: "base1",
	}

	config2 := testConfiguration{
		baseURL: "base2",
	}

	underTest := hierarchicalConfigurationFormat{
		configurations: []Configuration{config1, config2},
	}

	assert.Equal(t, "base2", underTest.BaseURL())
}

func testOverrideHeaders(t *testing.T) {
	headerName := "someHeader"

	config1Headers := make(map[string][]string)
	config1Headers[headerName] = []string{"value1"}
	config1 := testConfiguration{
		headers: config1Headers,
	}

	config2Headers := make(map[string][]string)
	config2Headers[headerName] = []string{"value2"}
	config2 := testConfiguration{
		headers: config2Headers,
	}

	underTest := hierarchicalConfigurationFormat{
		configurations: []Configuration{config1, config2},
	}

	assert.Equal(t, "value2", underTest.Headers()[headerName][0])
}

func testOverrideVariables(t *testing.T) {
	variableName := "someVariable"

	config1Variables := make(map[string]string)
	config1Variables[variableName] = "value1"
	config1 := testConfiguration{
		variables: config1Variables,
	}

	config2Variables := make(map[string]string)
	config2Variables[variableName] = "value2"
	config2 := testConfiguration{
		variables: config2Variables,
	}

	underTest := hierarchicalConfigurationFormat{
		configurations: []Configuration{config1, config2},
	}

	assert.Equal(t, "value2", underTest.Variables()[variableName])
}
