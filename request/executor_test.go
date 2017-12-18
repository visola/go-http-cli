package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-http-cli/options"
)

func TestExecuteRequest(t *testing.T) {
	options := options.RequestOptions{
		URL: "https://www.google.com",
	}

	response, err := ExecuteRequest(options)

	assert.Nil(t, err, "Should execute request correctly")

	if err != nil {
		assert.Equal(t, response.Status, 200)
	}
}
