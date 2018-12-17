package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/visola/go-http-cli/model"
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
	Output  string
	Path    string
	Query   map[string]model.ArrayOrString
}

func checkExpected(spec *Spec, stdOut string) error {
	errorMessage := checkMethod(spec)
	errorMessage += checkHeaders(spec)
	errorMessage += checkPath(spec)
	errorMessage += checkBody(spec)
	errorMessage += checkOutput(spec.Expected.Output, stdOut)
	errorMessage += checkQueryParams(spec)

	if errorMessage != "" {
		return errors.New(errorMessage)
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

func checkOutput(expected string, actual string) string {
	expected = strings.TrimSpace(expected)
	expected = replaceVariablesInArray(expected)[0]

	if expected == "" {
		return ""
	}

	actual = strings.TrimSpace(actual)

	expectedSplit := strings.Split(strings.Replace(expected, "\r", "\n", -1), "\n")
	actualSplit := strings.Split(strings.Replace(actual, "\r", "\n", -1), "\n")

	if len(expectedSplit) != len(actualSplit) {
		return "Output has different number of lines (" + strconv.Itoa(len(expectedSplit)) +
			" vs " + strconv.Itoa(len(actualSplit)) + "):\n  Expected: \n--- Start\n" + expected +
			"\n---End\n    Actual: \n---Start\n" + actual + "\n---End\n"
	}

	linesFailed := make([]int, 0)
	linesNotFound := make([]int, 0)

	// Keeps track of lines in the actual that some other expected line already matched
	accountedLines := make(map[int]bool)

	for i, expectedLine := range expectedSplit {
		// Ignore the line
		if strings.HasPrefix(expectedLine, "#I# ") {
			expectedSplit[i] = expectedLine[4:] + " [IGNORED]"
			continue
		}

		// Unordered line, needs to be in the response
		if strings.HasPrefix(expectedLine, "#U# ") {
			expectedLine = expectedLine[4:]
			expectedSplit[i] = expectedLine
			if !isLinePresent(expectedLine, actualSplit, accountedLines) {
				linesNotFound = append(linesNotFound, i+1)
			}
		} else {
			accountedLines[i] = true
			if expectedLine != actualSplit[i] {
				linesFailed = append(linesFailed, i+1)
			}
		}
	}

	if len(linesFailed) > 0 || len(linesNotFound) > 0 {
		for _, lineFailed := range linesFailed {
			expectedSplit[lineFailed-1] = strings.TrimSpace(expectedSplit[lineFailed-1]) + " [FAILED]"
		}

		for _, lineNotFound := range linesNotFound {
			expectedSplit[lineNotFound-1] = expectedSplit[lineNotFound-1] + " [NOT FOUND]"
		}

		result := "Output doesn't match expected, failed lines " + fmt.Sprint(linesFailed)
		result += ", lines not found " + fmt.Sprint(linesNotFound) + ":\n"
		result += "Expected: --- Start\n" + addLineNumbers(expectedSplit) + "\n---End\n"
		result += "Actual: ---Start\n" + addLineNumbers(actualSplit) + "\n---End\n"
		return result
	}
	return ""
}

func checkQueryParams(spec *Spec) string {
	result := ""
	for key, expectedValues := range spec.Expected.Query {
		actualValues, present := lastRequest.Query[key]
		if !present {
			result += "Expected query parameter '" + key + "' but not found"
			continue
		}

		sort.Strings(expectedValues)
		sort.Strings(actualValues)
		if !reflect.DeepEqual([]string(expectedValues), []string(actualValues)) {
			result += "Query values don't match expected for key '" + key + "':\n"
			result += "  Expected: " + strings.Join(expectedValues, ",") + "\n"
			result += "    Actual: " + strings.Join(actualValues, ",") + "\n"

		}
	}
	return result
}

func addLineNumbers(lines []string) string {
	result := ""
	for i, line := range lines {
		result += string(i+1) + ". " + line + "\n"
	}
	return result
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

	mkdirErr := os.MkdirAll(spec.ProfileDir, 0777)
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

func isLinePresent(expectedLine string, actualLines []string, accountedLines map[int]bool) bool {
	for j, actualLine := range actualLines {
		if _, isAccountedFor := accountedLines[j]; isAccountedFor {
			continue
		}

		if actualLine == expectedLine {
			accountedLines[j] = true
			return true
		}
	}

	return false
}

func loadSpec(pathToSpecFile string) (*Spec, error) {
	data, readErr := ioutil.ReadFile(pathToSpecFile)
	if readErr != nil {
		return nil, readErr
	}

	loadedSpec := new(Spec)
	unmarshalErr := yaml.Unmarshal(data, loadedSpec)

	relPath, _ := filepath.Rel(specsFolder, pathToSpecFile)
	loadedSpec.ProfileDir = path.Join(".", relPath)

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
	exitCode, stdOut, stdErr, execError := executeCommand(loadedSpec.Command[0], replaceVariablesInArray(loadedSpec.Command[1:]...))
	if execError != nil {
		if stdErr != "" {
			return fmt.Errorf("%s\n-- Standard Error --\n%s", execError.Error(), stdErr)
		}
		return execError
	}

	if exitCode != 0 {
		return fmt.Errorf("Exit code wasn't 0: %d", exitCode)
	}

	return checkExpected(loadedSpec, stdOut)
}

func replaceVariablesInArray(arrayIn ...string) []string {
	result := make([]string, len(arrayIn))
	for i, val := range arrayIn {
		result[i] = variables.ReplaceVariables(val, getContext())
	}
	return result
}
