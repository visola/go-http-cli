package request

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequest(t *testing.T) {
	t.Run("Build request correctly", testBuildsRequestCorrectly)
}

func testBuildsRequestCorrectly(t *testing.T) {
	request := Request{
		Body:   `{"username":"John Doe"}`,
		Method: http.MethodPost,
		QueryParams: map[string][]string{
			"auth": {"4312763812&*&%&$%!^@#+123"},
		},
		URL: "http://www.someserver.com/1234/employee",
	}

	httpReq, reqErr := BuildRequest(request)

	assert.Nil(t, reqErr, "Should create request")
	if reqErr != nil {
		panic(reqErr)
	}

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee?auth=4312763812%26%2A%26%25%26%24%25%21%5E%40%23%2B123 HTTP/1.1\r\nHost: www.someserver.com\r\n\r\n{\"username\":\"John Doe\"}",
		string(dump),
		"Should generate the expected dump",
	)
}
