package integrationtests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/visola/go-http-cli/pkg/profile"

	"github.com/stretchr/testify/assert"
)

// ExecuteCommand executes an external command and wait for the command to finish executing
func ExecuteCommand(cmd string, args ...string) (int, string, string, error) {
	executionDir := os.Getenv("EXECUTION_DIR")
	command := exec.Command(cmd, args...)
	command.Dir = executionDir

	profilesDir, _ := profile.GetProfilesDir()
	env := []string{"PATH=" + executionDir, "GO_HTTP_PROFILES=" + profilesDir}
	command.Env = env

	var outbuf, errbuf bytes.Buffer
	command.Stdout = &outbuf
	command.Stderr = &errbuf

	execErr := command.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	if execErr != nil {
		if exitError, ok := execErr.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			execErr = fmt.Errorf("Error while executing command.\n%s\nstdout:\n%s\nstderr:\n%s", execErr.Error(), stdout, stderr)
			return ws.ExitStatus(), stdout, stderr, execErr
		}
		return -1, stdout, stderr, execErr
	}

	ws := command.ProcessState.Sys().(syscall.WaitStatus)
	exitCode := ws.ExitStatus()
	return exitCode, stdout, stderr, nil
}

// KillDaemon runs the kill daemon command
func KillDaemon() {
	ExecuteCommand("go-http-daemon", "--kill")
}

// RunHTTP will run the http CLI with the specified arguments, ensure that it finished correctly
// and return the output
func RunHTTP(t *testing.T, args ...string) string {
	_, output, _, err := ExecuteCommand("./http", args...)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	return output
}
