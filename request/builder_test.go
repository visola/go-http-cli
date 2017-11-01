package request

import (
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visola/go-http-cli/config"
)

func TestBuildRequest(t *testing.T) {
	config := config.TestConfiguration{
		TestBaseURL: "http://www.someserver.com/",
		TestBody:    "Some data to be sent.",
		TestHeaders: map[string][]string{
			"Content-Type": {"application/json"},
		},
		TestMethod: "POST",
		TestURL:    "/path/with/${someVariable}/something",
		TestVariables: map[string]string{
			"someVariable": "someValue",
		},
	}

	req, reqErr := BuildRequest(config)
	assert.Nil(t, reqErr, "Should create request")

	headers := make([]string, 0, len(req.Header))
	for k := range req.Header {
		headers = append(headers, k)
	}

	assert.Equal(t, 1, len(headers), "Should have one header set")
	assert.Equal(t, headers[0], "Content-Type")

	assert.Equal(t, "/path/with/someValue/something", req.URL.EscapedPath(), "Should replace variables in the path")

	dump, dumpErr := httputil.DumpRequest(req, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /path/with/someValue/something HTTP/1.1\r\nHost: www.someserver.com\r\nContent-Type: application/json\r\n\r\nSome data to be sent.",
		string(dump),
		"Should generate the expected dump",
	)
}
