package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
			start := time.Now().UnixNano()

			loadedTestCase, loadErr := loadTestCase(pathToFile)
			if loadErr != nil {
				total := (time.Now().UnixNano() - start) / int64(time.Millisecond)
				errorColor.Println("Failed to load test case (" + strconv.FormatInt(total, 10) + "ms): " + pathToFile + "\n" + loadErr.Error() + "\n")
				return nil
			}

			runErr := loadedTestCase.run()
			total := (time.Now().UnixNano() - start) / int64(time.Millisecond)
			if runErr != nil {
				errorColor.Println("Failed (" + strconv.FormatInt(total, 10) + "ms): " + pathToFile + "\n" + runErr.Error() + "\n")
				return nil
			}
			successColor.Printf("Passed (%dms): %s\n", total, pathToFile)
		}

		return nil
	})

	if walkErr != nil {
		panic(walkErr)
	}
}
