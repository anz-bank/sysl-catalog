package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

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
	var output string
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

	README.Write([]byte("| Service | Method |\n| - |:-:|\n"))

	for _, app := range m.Apps {
		appName := strings.Join(app.Name.GetPart(), "")
		fs.MkdirAll(path.Join(output, appName), os.ModePerm)
		for _, endpoint := range app.Endpoints {
			outputFileName := path.Join(appName, appName + endpoint.Name + ".png")
			README.Write([]byte(fmt.Sprintf("[%s](%s) | [%s](%s) \n", appName, appName, endpoint.Name, outputFileName)))
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
