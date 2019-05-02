package integration

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// WrapForIntegrationTest wraps a testing function with all the required pieces for an integration test
func WrapForIntegrationTest(toWrap func(*testing.T)) func(*testing.T) {
	os.Setenv("GO_HTTP_PROFILES", path.Join(os.Getenv("EXECUTION_DIR"), getTestName(toWrap)))
	return WrapWithKillDamon(WrapWithTestServer(toWrap))
}

// WrapWithKillDamon executes a test after calling KillDaemon
func WrapWithKillDamon(toWrap func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		KillDaemon()
		toWrap(t)
	}
}

// WrapWithTestServer initializes the test server and make sure it will tear down correctly after
func WrapWithTestServer(toWrap func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		startTestServer()
		defer testServer.Close()

		toWrap(t)
	}
}

func getTestName(testFunc func(*testing.T)) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(testFunc).Pointer()).Name()
	split := strings.Split(fullName, ".")
	return split[len(split)-1]
}
