package request

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleParseURL() {
	base := "http://localhost:3000/api/v1"
	url := "/{companyId}/contacts"
	context := map[string]string{
		"companyId": "123456",
	}
	fmt.Println(ParseURL(base, url, context))
	// Output: http://localhost:3000/api/v1/123456/contacts
}

func TestParseURL(t *testing.T) {
	t.Run("Overrides base URL if full URL is passed in.", testOverridesBaseURL)
}

func testOverridesBaseURL(t *testing.T) {
	base := "http://localhost:3000/api/v1"
	url := "https://localhost:3000/some/where/else"

	result := ParseURL(base, url, nil)

	assert.Equal(t, url, result, "Should override base URL")
}
