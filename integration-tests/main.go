package main

import (
	"fmt"
	"io/ioutil"
)

const specsFolder = "specs"

func main() {
	startTestServer()
	defer testServer.Close()
	fmt.Printf("Test Server is open at: %s\n", testServer.URL)

	files, filesErr := ioutil.ReadDir(specsFolder)
	if filesErr != nil {
		panic(filesErr)
	}

	for _, specFile := range files {
		fmt.Printf("Running %s\n", specFile.Name())
		runErr := runSpec(specFile)
		if runErr != nil {
			panic(runErr)
		}
	}
}
