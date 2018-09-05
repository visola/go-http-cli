package request

import (
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequest(t *testing.T) {
	t.Run("Build request works", testRequestBuilder)
	t.Run("Build GET request with values", testGetRequestWithValues)
	t.Run("Build POST with JSON from values", testPostRequestWithJSONFromValues)
	t.Run("Build POST with URL Encoded from values", testPostRequestWithURLEncodedFromValues)
}

func testGetRequestWithValues(t *testing.T) {
	request := Request{
		URL: "http://www.someserver.com/{companyId}/employee",
		Values: map[string][]string{
			"name": {"{username}"},
		},
	}

	variables := map[string]string{
		"companyId": "1234",
		"username":  "John Doe",
	}

	httpReq, reqErr := BuildRequest(request, nil, variables)

	assert.Nil(t, reqErr, "Should create request")
	if reqErr != nil {
		panic(reqErr)
	}

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"GET /1234/employee?name=John+Doe HTTP/1.1\r\nHost: www.someserver.com\r\n\r\n",
		string(dump),
		"Should generate the expected dump",
	)
}

func testPostRequestWithJSONFromValues(t *testing.T) {
	request := Request{
		Headers: map[string][]string{
			"Company-Id": {"{companyId}"},
		},
		Method: "PUT",
		URL:    "http://www.someserver.com/{companyId}/employee",
		Values: map[string][]string{
			"companyId": {"{companyId}"},
			"name":      {"{username}"},
		},
	}

	variables := map[string]string{
		"companyId": "1234",
		"username":  "John Doe",
	}

	httpReq, reqErr := BuildRequest(request, nil, variables)
	assert.Nil(t, reqErr, "Should create request")

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"PUT /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nCompany-Id: 1234\r\nContent-Type: application/json\r\n\r\n{\"companyId\":\"1234\",\"name\":\"John Doe\"}",
		string(dump),
		"Should generate the expected dump",
	)
}

func testPostRequestWithURLEncodedFromValues(t *testing.T) {
	request := Request{
		Headers: map[string][]string{
			"Company-Id":   {"{companyId}"},
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		URL: "http://www.someserver.com/{companyId}/employee",
		Values: map[string][]string{
			"companyId": {"{companyId}"},
			"name":      {"{username}"},
		},
	}

	variables := map[string]string{
		"companyId": "1234",
		"username":  "John Doe",
	}

	httpReq, reqErr := BuildRequest(request, nil, variables)
	assert.Nil(t, reqErr, "Should create request")

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nCompany-Id: 1234\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\ncompanyId=1234&name=John+Doe",
		string(dump),
		"Should generate the expected dump",
	)
}

func testRequestBuilder(t *testing.T) {
	request := Request{
		Body: `{"name":"{username},"companyId":{companyId}}`,
		Headers: map[string][]string{
			"Company-Id":   {"{companyId}"},
			"Content-Type": {"application/json"},
		},
		URL: "http://www.someserver.com/{companyId}/employee",
	}

	variables := map[string]string{
		"companyId": "1234",
		"username":  "John Doe",
	}

	httpReq, reqErr := BuildRequest(request, nil, variables)
	assert.Nil(t, reqErr, "Should create request")

	dump, dumpErr := httputil.DumpRequest(httpReq, true)
	assert.Nil(t, dumpErr, "Dump should work")
	assert.Equal(
		t,
		"POST /1234/employee HTTP/1.1\r\nHost: www.someserver.com\r\nCompany-Id: 1234\r\nContent-Type: application/json\r\n\r\n{\"name\":\"John Doe,\"companyId\":1234}",
		string(dump),
		"Should generate the expected dump",
	)
}
