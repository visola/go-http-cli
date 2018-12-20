package main

import (
	"io/ioutil"
	"path"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// TestCase represents each test case file, it can be just a spec, or it can contain multiple specs
type TestCase struct {
	Command    []string
	Expected   Expected
	ProfileDir string
	Profiles   map[string]string
	ReplyWith  *ReplyWith

	Specs []*Spec
}

func loadTestCase(pathToSpecFile string) (*TestCase, error) {
	data, readErr := ioutil.ReadFile(pathToSpecFile)
	if readErr != nil {
		return nil, readErr
	}

	loadedTestCase := new(TestCase)
	unmarshalErr := yaml.Unmarshal(data, loadedTestCase)

	relPath, _ := filepath.Rel(specsFolder, pathToSpecFile)
	loadedTestCase.ProfileDir = path.Join(".", relPath)

	return loadedTestCase, unmarshalErr
}

func (testCase *TestCase) run() error {
	executeCommand("go-http-daemon", []string{"--kill"})

	for _, loadedSpec := range testCase.toSpecs() {
		runErr := runSpec(loadedSpec)
		if runErr != nil {
			return runErr
		}
	}

	return nil
}

func (testCase *TestCase) toSpecs() []*Spec {
	result := make([]*Spec, 0)
	if len(testCase.Specs) == 0 {
		result = append(result, &Spec{
			Command:    testCase.Command,
			Expected:   testCase.Expected,
			ProfileDir: testCase.ProfileDir,
			Profiles:   testCase.Profiles,
		})
	} else {
		for _, spec := range testCase.Specs {
			result = append(result, spec)
		}
	}

	return result
}
