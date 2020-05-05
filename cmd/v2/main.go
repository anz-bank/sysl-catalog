package main

import (
	"os"

	"github.com/anz-bank/sysl/pkg/parse"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	filename := "tests/simple.sysl"
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	catalog.Run(m, plantumlService, fs)
}
