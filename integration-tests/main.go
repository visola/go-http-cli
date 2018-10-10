package main

import (
	"io/ioutil"

	"github.com/fatih/color"
)

const specsFolder = "../specs"

func main() {
	startTestServer()
	defer testServer.Close()

	files, filesErr := ioutil.ReadDir(specsFolder)
	if filesErr != nil {
		panic(filesErr)
	}

	errorColor := color.New(color.FgRed)
	successColor := color.New(color.FgGreen)
	for _, specFile := range files {
		runErr := runSpec(specFile)
		if runErr != nil {
			errorColor.Printf("Error while running spec: %s\n%s", specFile.Name(), runErr.Error())
			continue
		}
		successColor.Printf("Passed: %s", specFile.Name())
	}
}
