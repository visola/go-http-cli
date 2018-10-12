package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/visola/variables/variables"
	yaml "gopkg.in/yaml.v2"
)

// Spec is the struct that represents a test case
type Spec struct {
	Command  []string
	Expected Expected
}

// Expected is the struct that stores an expected result
type Expected struct {
	Method   string
	Response string
}

func checkExpected(spec *Spec) error {
	errorMessage := ""
	if spec.Expected.Method != lastRequest.Method {
		errorMessage += fmt.Sprintf("  - Unexpected HTTP Method: \n    Expected: %s\n      Actual: %s", spec.Expected.Method, lastRequest.Method)
	}

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

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
			return ws.ExitStatus(), stdout, stderr, execErr
		}
		return -1, stdout, stderr, execErr
	}

	ws := command.ProcessState.Sys().(syscall.WaitStatus)
	exitCode := ws.ExitStatus()
	return exitCode, stdout, stderr, nil
}

func loadSpec(pathToSpecFile string) (*Spec, error) {
	data, readErr := ioutil.ReadFile(pathToSpecFile)
	if readErr != nil {
		return nil, readErr
	}

	loadedSpec := new(Spec)
	unmarshalErr := yaml.Unmarshal(data, loadedSpec)
	return loadedSpec, unmarshalErr
}

func runSpec(pathToSpecFile string) error {
	loadedSpec, loadErr := loadSpec(pathToSpecFile)
	if loadErr != nil {
		return loadErr
	}

	exitCode, _, stdErr, execError := executeCommand(loadedSpec.Command[0], replaceVariablesInArray(loadedSpec.Command[1:]))
	if execError != nil {
		if stdErr != "" {
			return fmt.Errorf("%s\n-- Standard Error --\n%s", execError.Error(), stdErr)
		}
		return execError
	}

	if exitCode != 0 {
		return fmt.Errorf("Exit code wasn't 0: %d", exitCode)
	}

	return checkExpected(loadedSpec)
}

func replaceVariablesInArray(arrayIn []string) []string {
	result := make([]string, len(arrayIn))
	for i, val := range arrayIn {
		result[i] = variables.ReplaceVariables(val, getContext())
	}
	return result
}
