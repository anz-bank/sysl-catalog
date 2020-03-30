package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/anz-bank/sysl-catalog/pkg/templategeneration"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func main() {
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	var output string
	flag.StringVar(&output, "o", "./", "Output directory of documentation")
	flag.Parse()
	filename := flag.Arg(0)
	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}

	// fs = afero.NewMemMapFs()
	fs.MkdirAll("docs", os.ModePerm)
	file, _ := fs.Create("docs/index.html")
	file.Write([]byte("Hello, world"))
	file.Close()
	templategeneration.NewProject(filename, "docs", plantumlService, logrus.New(), fs, m).
		ExecuteTemplateAndDiagrams()

	httpFs := afero.NewHttpFs(fs)
	fileserver := http.FileServer(httpFs.Dir("/docs/"))
	http.Handle("/", fileserver)
	http.ListenAndServe(":8080", fileserver)
}
