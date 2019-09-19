package integration

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// WithTempFile creates a temporary file, write content to it and then call the wrapped function
func WithTempFile(t *testing.T, content string, toWrap func(*os.File)) {
	tempFile, err := ioutil.TempFile("", "script.js")
	if err != nil {
		t.Fatalf("Error while creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error while writing content to temp file: %s", err)
	}

	toWrap(tempFile)
}

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
