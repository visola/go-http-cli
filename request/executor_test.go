package request

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteRequest(t *testing.T) {
	t.Run("Request httpbin should return 200", testBasicGet)
	t.Run("Request httpbin redirect should follow redirect", testFollowsRedirect)
}

func testBasicGet(t *testing.T) {
	request := Request{
		Method: http.MethodGet,
		URL:    "https://httpbin.org/",
	}

	executedRequestResponses, err := ExecuteRequest(request, nil, nil)

	assert.Nil(t, err, "Should execute request correctly")

	if err == nil {
		assert.Equal(t, 1, len(executedRequestResponses))
		assert.Equal(t, http.StatusOK, executedRequestResponses[0].Response.StatusCode)
	}
}

func testFollowsRedirect(t *testing.T) {
	request := Request{
		URL: "https://httpbin.org/redirect/1",
	}

	executedRequestResponses, err := ExecuteRequest(request, nil, nil)

	assert.Nil(t, err, "Should execute request correctly")

	if err == nil {
		assert.Equal(t, 2, len(executedRequestResponses), "Should have executed 2 requests")
		assert.Equal(t, http.StatusFound, executedRequestResponses[0].Response.StatusCode, "First response should be 302")
		assert.Equal(t, http.StatusOK, executedRequestResponses[1].Response.StatusCode, "Second response should be 200")
	}
}
