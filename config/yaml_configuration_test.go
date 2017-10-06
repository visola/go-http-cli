package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFrom(t *testing.T) {
	headerName := "Content-type"
	headerValue := "application/json"

	simpleHeaderYaml := "headers:\n  " + headerName + ":\n    - " + headerValue
	tmpFile, err := ioutil.TempFile("", "simple_header.yml")

	if err != nil {
		panic(err)
	}

	tmpFile.WriteString(simpleHeaderYaml)

	parsedConfiguration, err := readFrom(tmpFile.Name())

	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, len(parsedConfiguration.Headers()), 1, "Should have parsed one header")
	assert.Equal(t, parsedConfiguration.Headers()[headerName], []string{headerValue}, "Should parse header value correctly")
}
