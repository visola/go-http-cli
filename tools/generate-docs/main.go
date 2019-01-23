package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DocumentPart represents one part of the final document
type DocumentPart struct {
	Title   string
	Content string
}

func main() {
	f, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, part := range collectAllParts() {
		f.WriteString(part.Content)
		f.WriteString("\n")
	}
}

func collectAllParts() []DocumentPart {
	allParts := make([]DocumentPart, 0)
	filepath.Walk("docs/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		docPart, docPartErr := getDocumentPart(path)
		if docPartErr != nil {
			return docPartErr
		}

		allParts = append(allParts, docPart)

		return nil
	})

	return allParts
}

func getDocumentPart(path string) (docPart DocumentPart, err error) {
	content, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		err = readErr
		return
	}

	fileName := filepath.Base(path)
	docPart = DocumentPart{
		Content: string(content),
		Title:   getTitleFromFileName(fileName),
	}

	return
}

func getTitleFromFileName(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	return strings.Split(fileName[0:len(fileName)-len(fileExtension)], "_")[1]
}
