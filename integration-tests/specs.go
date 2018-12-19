package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/visola/variables/variables"
)

// Spec is the struct that represents a test case
type Spec struct {
	Command    []string
	Expected   Expected
	ProfileDir string
	Profiles   map[string]string
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

func runSpec(loadedSpec *Spec) error {
	os.Setenv("GO_HTTP_PROFILES", loadedSpec.ProfileDir)

	writeProfilesErr := createProfiles(loadedSpec)
	if writeProfilesErr != nil {
		return writeProfilesErr
	}

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
