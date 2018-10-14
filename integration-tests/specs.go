package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"syscall"

	"github.com/visola/variables/variables"
	yaml "gopkg.in/yaml.v2"
)

// Spec is the struct that represents a test case
type Spec struct {
	Command    []string
	Expected   Expected
	ProfileDir string
	Profiles   map[string]string
}

// Expected is the struct that stores an expected result
type Expected struct {
	Body    string
	Headers map[string][]string
	Method  string
	Path    string
}

func checkExpected(spec *Spec) error {
	errorMessage := checkMethod(spec)
	errorMessage += checkHeaders(spec)
	errorMessage += checkPath(spec)
	errorMessage += checkBody(spec)

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func checkBody(spec *Spec) string {
	errorMessage := ""
	if spec.Expected.Body != lastRequest.Body {
		errorMessage += "  - Expected body doesn't match:\n"
		errorMessage += fmt.Sprintf("Bodies:\nExpected:\n---\n%s\n---\nActual:\n--\n%s\n--\n", spec.Expected.Body, lastRequest.Body)
	}

	return errorMessage
}

func checkHeaders(spec *Spec) string {
	errorMessage := ""
	if len(spec.Expected.Headers) > 0 {
		headerError := false
		for expectedHeader, expectedValues := range spec.Expected.Headers {
			actualValues, headerExist := lastRequest.Headers[expectedHeader]
			if !headerExist {
				errorMessage += fmt.Sprintf(" - Expected header not found: %s\n", expectedHeader)
				headerError = true
				continue
			}

			sort.Strings(expectedValues)
			sort.Strings(actualValues)
			if !reflect.DeepEqual(expectedValues, actualValues) {
				headerError = true
				errorMessage += fmt.Sprintf(" - Header values do not match for header: %s\n  Expected: %s\n    Actual: %s\n", expectedHeader, expectedValues, actualValues)
			}
		}
		if headerError {
			errorMessage += fmt.Sprintf("Headers found:\n%s\n", lastRequest.Headers)
		}
	}

	return errorMessage
}

func checkMethod(spec *Spec) string {
	if spec.Expected.Method != "" && spec.Expected.Method != lastRequest.Method {
		return fmt.Sprintf("Unexpected HTTP Method: \n  Expected: %s\n    Actual: %s\n", spec.Expected.Method, lastRequest.Method)
	}
	return ""
}

func checkPath(spec *Spec) string {
	if spec.Expected.Path != "" && spec.Expected.Path != lastRequest.Path {
		return fmt.Sprintf("Unexpected Path: \n  Expected: %s\n    Actual: %s\n", spec.Expected.Path, lastRequest.Path)
	}
	return ""
}

func createProfiles(spec *Spec) error {
	if len(spec.Profiles) == 0 {
		return nil
	}

	mkdirErr := os.Mkdir(spec.ProfileDir, 0777)
	if mkdirErr != nil {
		return mkdirErr
	}

	for name, content := range spec.Profiles {
		profileFileName := fmt.Sprintf("%s.yaml", name)
		content = variables.ReplaceVariables(content, getContext())
		writeErr := ioutil.WriteFile(path.Join(spec.ProfileDir, profileFileName), []byte(content), 0777)
		if writeErr != nil {
			return writeErr
		}
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
			execErr = fmt.Errorf("Error while executing command.\n%s\nstdout:\n%s\nstderr:\n%s", execErr.Error(), stdout, stderr)
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

	_, specFile := filepath.Split(pathToSpecFile)
	loadedSpec.ProfileDir = path.Join(".", specFile)

	return loadedSpec, unmarshalErr
}

func runSpec(pathToSpecFile string) error {
	executeCommand("go-http-daemon", []string{"--kill"})

	loadedSpec, loadErr := loadSpec(pathToSpecFile)
	if loadErr != nil {
		return loadErr
	}

	writeProfilesErr := createProfiles(loadedSpec)
	if writeProfilesErr != nil {
		return writeProfilesErr
	}

	os.Setenv("GO_HTTP_PROFILES", loadedSpec.ProfileDir)
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
