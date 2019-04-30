package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DocumentationContext is what will be passed to the templates when generating the final README
type DocumentationContext struct {
	Parts []DocumentPart
}

// DocumentPart represents one part of the final document
type DocumentPart struct {
	Anchor  string
	Content string
	Path    string
	Title   string
}

func main() {
	f, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	context := DocumentationContext{
		Parts: collectAllParts(),
	}

	for _, part := range context.Parts {
		tmpl := template.Must(template.New(part.Path).Parse(part.Content))
		execErr := tmpl.Execute(f, context)
		if execErr != nil {
			panic(execErr)
		}
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

		expectedFileFormat, _ := regexp.MatchString(`\d+_[\w\s]+\.md`, info.Name())
		if !expectedFileFormat {
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
	title := getTitleFromFileName(fileName)
	anchor := strings.Replace(strings.ToLower(title), " ", "-", -1)

	docPart = DocumentPart{
		Anchor:  anchor,
		Content: string(content),
		Path:    path,
		Title:   title,
	}

	return
}

func getTitleFromFileName(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	return strings.Split(fileName[0:len(fileName)-len(fileExtension)], "_")[1]
}
