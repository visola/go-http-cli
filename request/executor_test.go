package request

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteRequest(t *testing.T) {
	t.Run("Request httpbin should return 200", testBasicGet)
	t.Run("Request httpbin redirect should follow redirect", testFollowsRedirect)
	t.Run("Should bail if max number of redirects happens", testMaxRedirects)
}

func testBasicGet(t *testing.T) {
	request := Request{
		Method: http.MethodGet,
		URL:    "https://httpbin.org/",
	}

	executedRequestResponses, err := ExecuteRequestLoop(ExecutionContext{
		Request: request,
	})

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

	executedRequestResponses, err := ExecuteRequestLoop(ExecutionContext{
		FollowLocation: true,
		Request:        request,
	})

	assert.Nil(t, err, "Should execute request correctly")

	if err == nil {
		assert.Equal(t, 2, len(executedRequestResponses), "Should have executed 2 requests")
		assert.Equal(t, http.StatusFound, executedRequestResponses[0].Response.StatusCode, "First response should be 302")
		assert.Equal(t, http.StatusOK, executedRequestResponses[1].Response.StatusCode, "Second response should be 200")
	}
}

func testMaxRedirects(t *testing.T) {
	const maxRedirectCount = 10

	request := Request{
		URL: fmt.Sprintf("https://httpbin.org/redirect/%d", maxRedirectCount+1),
	}

	executedRequestResponses, err := ExecuteRequestLoop(ExecutionContext{
		FollowLocation: true,
		MaxRedirect:    maxRedirectCount,
		Request:        request,
	})

	assert.NotNil(t, err, "Should return an error")

	// It should still return the requests that were executed and their responses
	assert.Equal(t, 11, len(executedRequestResponses), "Should have executed 11 requests")
}
