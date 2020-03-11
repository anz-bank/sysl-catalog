package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

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
	filename := os.Args[1]

	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	t, err := template.New("markdown").Parse(Index)
	README, err := fs.Create("README.md")
	if err != nil {
		panic(err)
	}
	err = t.Execute(README, m)

	for _, app := range m.Apps {
		appName := strings.Join(app.Name.GetPart(), "")
		for _, endpoint := range app.Endpoints {
			outputFileName := appName + endpoint.Name + ".png"
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

const Index = `| Service | Method |
| - |:-:|{{range $key, $val := .Apps}}{{range $endpointName, $endpoint := $val.Endpoints}}
|{{$key}}|{{$endpointName}}|{{end}}{{end}}
`
