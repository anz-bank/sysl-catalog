package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/spf13/afero"
)

func main() {
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	var output, packageName string
	flag.StringVar(&output, "o", "./", "Output directory of documentation")
	flag.Parse()
	filename := flag.Arg(0)
	fmt.Println(filename)
	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	README, err := fs.Create(output + "/README.md")

	if err != nil {
		panic(err)
	}

	README.Write([]byte("| Package |\n| - |\n"))
	packageReadmes := make(map[string]afero.File)
	var pacakgeREADME afero.File
	var ok bool
	for _, app := range m.Apps {
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		appName := strings.Join(app.Name.GetPart(), "")

		if attr := app.GetAttrs()["package"]; attr != nil {
			packageName = attr.GetS()
		} else {
			packageName = appName
		}
		fs.MkdirAll(path.Join(output, packageName), os.ModePerm)
		packageReadmeName := filepath.Join(output, packageName, "Readme.md")
		if pacakgeREADME, ok = packageReadmes[packageReadmeName]; !ok {
			pacakgeREADME, err = fs.Create(packageReadmeName)
			packageReadmes[packageReadmeName] = pacakgeREADME
			pacakgeREADME.Write([]byte("| Service | Method |\n| - |:-:|\n"))
			README.Write([]byte(fmt.Sprintf("[%s](%s) |\n", packageName, packageName)))
			if err != nil {
				panic(err)
			}
		} else {
			pacakgeREADME = packageReadmes[packageReadmeName]
		}

		for _, endpoint := range app.Endpoints {
			outputFileName := path.Join(output, packageName, appName+endpoint.Name+".svg")
			pacakgeREADME.Write([]byte(fmt.Sprintf("%s | [%s](/%s) \n", packageName, appName, outputFileName)))
			if err != nil {
				panic(err)
			}
			sequenceDiagram, err := CreateSequenceDiagram(m, fmt.Sprintf("%s <- %s", appName, endpoint.Name))
			if err != nil {
				panic(err)
			}
			diagrams.OutputPlantuml(outputFileName, plantumlService, sequenceDiagram, fs)
		}
	}
	if err != nil {
		panic(err)
	}
}
func CreateSequenceDiagram(m *sysl.Module, call string) (string, error) {
	l := &cmdutils.Labeler{}
	p := &sequencediagram.SequenceDiagParam{}
	p.Endpoints = []string{call}
	p.AppLabeler = l
	p.EndpointLabeler = l
	p.Title = call
	return sequencediagram.GenerateSequenceDiag(m, p, logrus.New())
}
