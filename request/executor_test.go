package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteRequest(t *testing.T) {
	request := Request{
		Method: "GET",
		URL:    "https://www.google.com",
	}

	executedRequestResponse, err := ExecuteRequest(request, nil, nil)

	assert.Nil(t, err, "Should execute request correctly")

	if err != nil {
		assert.Equal(t, executedRequestResponse.Response.Status, 200)
	}
}
