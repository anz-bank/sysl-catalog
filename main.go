package main

import (
	"flag"
	"fmt"
	"html/template"
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

type IndexMarkdown struct {
	PackageName    string
	PackageRelLink string
	App            []*AppMarkdown
}

const IndexMarkdownTemplate = `
| Package | Service Name | Method |
| - | - | - |{{range $PackxageName, $App := .}}{{range $Service := $App.App}}
[{{$App.PackageName}}]({{$App.PackageRelLink}})|{{$Service.ServiceName}}|[{{$Service.Method}}]({{$Service.RelLink}}) |{{end}}{{end}}
`

type AppMarkdown struct {
	ServiceName string
	Method      string
	Link        string
	RelLink     string
}

const AppMarkdownTemplate = `
[Back](../README.md)
| Service | Method |
| - |:-:|
{{range $App := .}}{{$App.ServiceName}}|({{$App.Method}})[{{$App.RelLink}}] |
{{end}}
`

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
	Index := make(map[string]*IndexMarkdown)

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
		packageRelLink := filepath.Join(packageName, "README.md")
		MarkdownApp := []*AppMarkdown{}
		fs.MkdirAll(path.Join(output, packageName), os.ModePerm)
		for _, endpoint := range app.Endpoints {
			outputFilePath := path.Join(output, packageName)
			outputFileName := path.Join(output, packageName, endpoint.Name+".svg")
			MarkdownApp = append(MarkdownApp, &AppMarkdown{
				ServiceName: appName,
				Method:      endpoint.Name,
				Link:        outputFilePath,
				RelLink:     endpoint.Name + ".svg",
			})
			CreateSequenceDiagramFile(
				m,
				fmt.Sprintf("%s <- %s", appName, endpoint.Name),
				outputFileName,
				plantumlService,
				fs)
		}
		if _, ok := Index[packageRelLink]; !ok {
			Index[packageRelLink] = &IndexMarkdown{PackageName: packageName, PackageRelLink: packageRelLink}
			Index[packageRelLink].App = []*AppMarkdown{}
		}

		Index[packageRelLink].App = append(Index[packageRelLink].App, MarkdownApp...)
		fmt.Println(Index[packageRelLink].App)
	}

	README, err := fs.Create(output + "/README.md")
	if err != nil {
		panic(err)
	}
	IndexTemplate, err := template.New("markdown").Parse(IndexMarkdownTemplate)
	err = IndexTemplate.Execute(README, Index)
	if err != nil {
		panic(err)
	}
	AppTemplate, err := template.New("markdown").Parse(AppMarkdownTemplate)
	if err != nil {
		panic(err)
	}
	for _, Apps := range Index {
		README, err := fs.Create(path.Join(output, Apps.PackageRelLink))
		if err != nil {
			panic(err)
		}
		fmt.Println("Creating", README.Name(), Apps)
		err = AppTemplate.Execute(README, Apps.App)
	}

}

//func GenerateMarkdown(index []IndexMarkdown, outputName string, fs){
//
//}

func CreateSequenceDiagramFile(m *sysl.Module, call, outputFileName, plantumlService string, fs afero.Fs) error {
	sequenceDiagram, err := CreateSequenceDiagram(m, call)
	if err != nil {
		panic(err)
	}
	return diagrams.OutputPlantuml(outputFileName, plantumlService, sequenceDiagram, fs)
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
