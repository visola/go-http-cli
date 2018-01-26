package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	t.Run("Overrides base URL if full URL is passed in.", testOverridesBaseURL)
}

func testOverridesBaseURL(t *testing.T) {
	base := "http://localhost:3000/api/v1"
	url := "https://localhost:3000/some/where/else"

	result := ParseURL(base, url)

	assert.Equal(t, url, result, "Should override base URL")
}
