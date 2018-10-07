package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

	var outbuf, errbuf bytes.Buffer
	command.Stdout = &outbuf
	command.Stderr = &errbuf

	execErr := command.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	ws := command.ProcessState.Sys().(syscall.WaitStatus)
	exitCode := ws.ExitStatus()

	if execErr != nil {
		return exitCode, stdout, stderr, execErr
	}

	return exitCode, stdout, stderr, nil
}

func loadSpec(specFile os.FileInfo) (*Spec, error) {
	data, readErr := ioutil.ReadFile(path.Join(specsFolder, specFile.Name()))
	if readErr != nil {
		return nil, readErr
	}

	loadedSpec := new(Spec)
	unmarshalErr := yaml.Unmarshal(data, loadedSpec)
	return loadedSpec, unmarshalErr
}

func runSpec(specFile os.FileInfo) error {
	loadedSpec, loadErr := loadSpec(specFile)
	if loadErr != nil {
		return loadErr
	}

	exitCode, _, _, execError := executeCommand(loadedSpec.Command[0], replaceVariablesInArray(loadedSpec.Command[1:]))
	if execError != nil {
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
