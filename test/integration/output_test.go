package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	t.Run("Replace variables on output", WrapForIntegrationTest(testVariablesGetReplacedOnOutput))
}

func testVariablesGetReplacedOnOutput(t *testing.T) {
	companyID := "1234"
	output := RunHTTP(
		t,
		"-V", "companyId="+companyID,
		testServer.URL+"/companies/{companyId}",
	)

	lines := strings.Split(output, "\n")

	expectedFirstLine := fmt.Sprintf("GET %s/companies/%s", testServer.URL, companyID)
	assert.Equal(t, expectedFirstLine, lines[0], "Should replace variables on output")
	assert.Equal(t, "200 OK 1.1", lines[2], "Third line should show status")
}
