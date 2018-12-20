package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/visola/go-http-cli/model"
)

// Expected is the struct that stores an expected result
type Expected struct {
	Body    string
	Headers map[string][]string
	Method  string
	Output  string
	Path    string
	Query   map[string]model.ArrayOrString
}

func addLineNumbers(lines []string) string {
	result := ""
	for i, line := range lines {
		result += string(i+1) + ". " + line + "\n"
	}
	return result
}

func checkBody(spec *Spec) string {
	errorMessage := ""
	if spec.Expected.Body != lastRequest.Body {
		errorMessage += "  - Expected body doesn't match:\n"
		errorMessage += fmt.Sprintf("Bodies:\nExpected:\n---\n%s\n---\nActual:\n--\n%s\n--\n", spec.Expected.Body, lastRequest.Body)
	}

	return errorMessage
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

func checkPath(spec *Spec) string {
	if spec.Expected.Path != "" && spec.Expected.Path != lastRequest.Path {
		return fmt.Sprintf("Unexpected Path: \n  Expected: %s\n    Actual: %s\n", spec.Expected.Path, lastRequest.Path)
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