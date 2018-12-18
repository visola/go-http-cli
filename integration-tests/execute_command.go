package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func executeCommand(cmd string, args []string) (int, string, string, error) {
	command := exec.Command(cmd, args...)
	command.Env = os.Environ()

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
