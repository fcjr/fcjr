package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	_ "embed"
)

//go:embed README.md.tmpl
var readmeTemplate string

type ReadmeData struct{}

func main() {

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Fprintln(os.Stderr, "No output path provided")
		os.Exit(1)
	}
	var outPath = os.Args[1]

	outFile, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	var readmeData ReadmeData

	err = renderReadmeToWriter(&readmeData, outFile)
	if err != nil {
		panic(err)
	}
}

func renderReadmeToWriter(readmeData *ReadmeData, writer io.Writer) error {
	readmeTemplate := template.Must(
		template.New("").Parse(readmeTemplate),
	)

	return readmeTemplate.ExecuteTemplate(writer, "", readmeData)
}
