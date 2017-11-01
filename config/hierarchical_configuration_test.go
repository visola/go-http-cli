package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHierarchicalConfiguration(t *testing.T) {
	t.Run("Overrides base url correctly", testOverrideBaseURL)
	t.Run("Overrides headers correctly", testOverrideHeaders)
	t.Run("Overrides variables correctly", testOverrideVariables)
}

func testOverrideBaseURL(t *testing.T) {
	config1 := TestConfiguration{
		TestBaseURL: "base1",
	}

	config2 := TestConfiguration{
		TestBaseURL: "base2",
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
	config1 := TestConfiguration{
		TestHeaders: config1Headers,
	}

	config2Headers := make(map[string][]string)
	config2Headers[headerName] = []string{"value2"}
	config2 := TestConfiguration{
		TestHeaders: config2Headers,
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
	config1 := TestConfiguration{
		TestVariables: config1Variables,
	}

	config2Variables := make(map[string]string)
	config2Variables[variableName] = "value2"
	config2 := TestConfiguration{
		TestVariables: config2Variables,
	}

	underTest := hierarchicalConfigurationFormat{
		configurations: []Configuration{config1, config2},
	}

	assert.Equal(t, "value2", underTest.Variables()[variableName])
}
