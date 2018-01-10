package request

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequest(t *testing.T) {
	request := Request{
		Body: "Some data to be sent.",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Method: http.MethodPost,
		URL:    "http://www.someserver.com/{companyId}/employee",
	}

	variables := map[string]string{
		"companyId": "1234",
	}

	httpReq, _, reqErr := BuildRequest(request, nil, variables)
	assert.Nil(t, reqErr, "Should create request")

	headers := make([]string, 0, len(httpReq.Header))
	for k := range httpReq.Header {
		headers = append(headers, k)
	}

	assert.Equal(t, 1, len(headers), "Should have one header set")
	assert.Equal(t, headers[0], "Content-Type", "Should be content type")

	assert.Equal(t, "/1234/employee", httpReq.URL.EscapedPath(), "Should set URL")

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nContent-Type: application/json\r\n\r\nSome data to be sent.",
		string(dump),
		"Should generate the expected dump",
	)
}
