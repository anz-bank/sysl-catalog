package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl-catalog/pkg/templategeneration"

	"github.com/anz-bank/sysl/pkg/parse"
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
	templategeneration.NewProject(filename, output, plantumlService, logrus.New(), fs, m).ExecuteTemplateAndDiagrams()
}
