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
	Arguments []string
	Expected  Expected
}

// Expected is the struct that stores an expected result
type Expected struct {
	Body   string
	Method string
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

	exitCode, output, errorOut, execError := executeCommand("http", replaceVariablesInArray(loadedSpec.Arguments))
	if execError != nil {
		return execError
	}

	fmt.Printf("Executed Spec, Exit code: %d\nOutput:\n%s\nError:\n%s\n", exitCode, output, errorOut)

	return nil
}

func replaceVariablesInArray(arrayIn []string) []string {
	result := make([]string, len(arrayIn))
	for i, val := range arrayIn {
		result[i] = variables.ReplaceVariables(val, getContext())
	}
	return result
}
