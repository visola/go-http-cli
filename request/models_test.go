package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMashalUnmarshalExecutedRequestResponse(t *testing.T) {
	req := Request{
		Body: "Hello world!",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Method: "GET",
		URL:    "http://www.google.com",
	}

	resp := Response{
		Body: "Good Bye World!",
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Protocol:   "1.1",
		Status:     "OK",
		StatusCode: 200,
	}

	pair := ExecutedRequestResponse{
		Request:  req,
		Response: resp,
	}

	b, err := json.Marshal(pair)
	assert.Nil(t, err, "Should marshal correctly")

	var newPair ExecutedRequestResponse
	err = json.Unmarshal(b, &newPair)

	assert.Nil(t, err, "Should unmarshal correctly")
	assert.Equal(t, pair, newPair)
}
