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
	EndPoints      []*AppMarkdown
}

const IndexMarkdownTemplate = `
| Package | Service Name | EndpointName |
| - | - | - |{{range $PackageName, $App := .}}{{range $Service := $App.EndPoints}}
[{{$App.PackageName}}]({{$App.PackageRelLink}})|{{$Service.AppName}}|[{{$Service.EndpointName}}]({{$Service.Package}}/{{$Service.EndpointName}}.svg.md) |{{end}}{{end}}
`

type AppMarkdown struct {
	Package      string
	AppName      string
	EndpointName string
}

const AppMarkdownTemplate = `
[Back](../README.md)
| Service | EndpointName |
| - |:-:|
{{range $EndPoints := .}}{{$EndPoints.AppName}}|[{{$EndPoints.EndpointName}}]({{$EndPoints.EndpointName}}.svg.md) |
{{end}}
`
const embededSvgTemplate = `
[Back](README.md)

![alt text]({{.EndpointName}}.svg)

`
const ext = ".svg"
const generatedReadmeName = "README.md"

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
	Index := ConvertToMarkdownObject(m, output, plantumlService, fs)
	GenerateMarkdown(Index, output, fs)
}
func ConvertToMarkdownObject(m *sysl.Module, output, plantumlService string, fs afero.Fs) map[string]*IndexMarkdown {
	Index := make(map[string]*IndexMarkdown)
	var packageName string
	for _, app := range m.Apps {
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		appName := strings.Join(app.Name.GetPart(), "")
		packageName = appName
		if attr := app.GetAttrs()["package"]; attr != nil {
			packageName = attr.GetS()
		}
		packageRelLink := filepath.Join(packageName, generatedReadmeName)
		MarkdownApp := make([]*AppMarkdown, 0, len(app.Endpoints))
		fs.MkdirAll(path.Join(output, packageName), os.ModePerm)
		for _, endpoint := range app.Endpoints {
			outputFileName := path.Join(output, packageName, endpoint.Name+ext)
			MarkdownApp = append(MarkdownApp, &AppMarkdown{
				Package:      packageName,
				AppName:      appName,
				EndpointName: endpoint.Name,
			})

			err := CreateSequenceDiagramFile(
				m,
				fmt.Sprintf("%s <- %s", appName, endpoint.Name),
				outputFileName,
				plantumlService,
				fs)
			if err != nil {
				panic(err)
			}
		}
		if _, ok := Index[packageRelLink]; !ok {
			Index[packageRelLink] = &IndexMarkdown{PackageName: packageName, PackageRelLink: packageRelLink}
			Index[packageRelLink].EndPoints = []*AppMarkdown{}
		}
		Index[packageRelLink].EndPoints = append(Index[packageRelLink].EndPoints, MarkdownApp...)
	}
	return Index

}

func GenerateMarkdown(Index map[string]*IndexMarkdown, output string, fs afero.Fs) {
	IndexTemplate, err := template.New("markdown").Parse(IndexMarkdownTemplate)
	if err != nil {
		panic(err)
	}
	AppTemplate, err := template.New("markdown").Parse(AppMarkdownTemplate)
	if err != nil {
		panic(err)
	}
	EmbededSvgTemplate, err := template.New("markdown").Parse(embededSvgTemplate)
	if err != nil {
		panic(err)
	}
	README, err := fs.Create(output + "/README.md")
	fmt.Println("Creating", README.Name())
	if err != nil {
		panic(err)
	}
	if err := IndexTemplate.Execute(README, Index); err != nil {
		panic(err)
	}
	for _, Apps := range Index {
		README, err := fs.Create(path.Join(output, Apps.PackageRelLink))
		if err != nil {
			panic(err)
		}
		fmt.Println("Creating", README.Name())
		if err := AppTemplate.Execute(README, Apps.EndPoints); err != nil {
			panic(err)
		}
		for _, Endpoint := range Apps.EndPoints {
			embededSvg, err := fs.Create(path.Join(output, Endpoint.Package, Endpoint.EndpointName+ext+".md"))
			if err != nil {
				panic(err)
			}
			if err := EmbededSvgTemplate.Execute(embededSvg, Endpoint); err != nil {
				panic(err)
			}
		}
	}
}

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
