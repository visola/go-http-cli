package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const specsFolder = "../specs"

func main() {
	startTestServer()
	defer testServer.Close()

	errorColor := color.New(color.FgRed)
	successColor := color.New(color.FgGreen)
	walkErr := filepath.Walk(specsFolder, func(pathToFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), "yaml") || strings.HasSuffix(info.Name(), "yml") {
			runErr := runSpec(pathToFile)
			if runErr != nil {
				errorColor.Printf("Failed: %s\n%s\n", pathToFile, runErr.Error())
				return nil
			}
			successColor.Printf("Passed: %s\n", pathToFile)
		}

		return nil
	})

	if walkErr != nil {
		panic(walkErr)
	}
}
