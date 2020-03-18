package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/spf13/afero"
)

type IndexMarkdown struct {
	PackageName     string
	PackageRelLink  string
	EndPoints       []*AppMarkdown
	Module          *sysl.Module
	OutputDir       string
	Fs              afero.Fs
	PlantumlService string
}

type AppMarkdown struct {
	Package      string
	AppName      string
	EndpointName string
}

const ext = ".svg"

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
	p := NewProject(m.String(), output, plantumlService, IndexMarkdownTemplate, AppMarkdownTemplate, embededSvgTemplate, fs, m)
	p.GenerateMarkdownAndDiagrams()
}

func (p *Project) GenerateMarkdownAndDiagrams() {
	if err := GenerateMarkdown2(p.Output, p.OutputFile, p, p.ProjectTempl, p.Fs); err != nil {
		panic(err)
	}
	for _, key := range alphabeticalPackage(p.rows) {
		row := p.rows[key]
		if err := GenerateMarkdown2(row.OutputDir, row.OutputFile, row, row.Templ, p.Fs); err != nil {
			panic(err)
		}
		for _, sd := range row.SequenceDiagrams {
			if err := GenerateDiagramAndMarkdown(sd); err != nil {
				panic(err)
			}
		}
		for _, int := range row.IntegrationDiagrams {
			if err := GenerateDiagramAndMarkdown(int); err != nil {
				panic(err)
			}
		}
		for _, data := range row.DataModelDiagrams {
			if err := GenerateDiagramAndMarkdown(data); err != nil {
				panic(err)
			}
		}
	}
}

func GenerateDiagramAndMarkdown(sd *Diagram) error {
	fmt.Println(sd.OutputDirectory, sd.Name+".md")
	if err := GenerateMarkdown2(sd.OutputDirectory, sd.Name+".md", sd, sd.Templ, sd.Parent.Parent.Fs); err != nil {
		return err
	}
	outputFileName := path.Join(sd.OutputDirectory, sd.Name+ext)
	return diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.Diagram, sd.Parent.Parent.Fs)
}

func GenerateMarkdown2(outputdir, outputName string, object interface{}, t *template.Template, fs afero.Fs) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, object); err != nil {
		return err
	}
	fs.MkdirAll(outputdir, os.ModePerm)
	return afero.WriteFile(fs, filepath.Join(outputdir, outputName), buf.Bytes(), os.ModePerm)
}

//func GenerateMarkdown(project *Project) {
//
//	templates, err := LoadMarkdownTemplates(IndexMarkdownTemplate, AppMarkdownTemplate, embededSvgTemplate)
//	if err != nil {
//		panic(err)
//	}
//	indexTemplate, projectTemplate, _ := templates[0], templates[1], templates[2]
//	README, err := project.Fs.Create(project.Output + "/README.md")
//	if err != nil {
//		panic(err)
//	}
//	if err := indexTemplate.Execute(README, project); err != nil {
//		panic(err)
//	}
//
//	for _, key := range alphabeticalPackage(project.rows) {
//		row := project.rows[key]
//		projectReadme, err := project.Fs.Create(path.Join(row.OutputDir, "README.md"))
//		if err != nil {
//			panic(err)
//		}
//		if err := projectTemplate.Execute(projectReadme, project); err != nil {
//			panic(err)
//		}
//		for _, sd := range row.SequenceDiagrams {
//			outputFileName := path.Join(sd.OutputDirectory, sd.Name+ext)
//			project.Fs.MkdirAll(sd.OutputDirectory, os.ModePerm)
//			if err := diagrams.OutputPlantuml(outputFileName, project.PlantumlService, sd.Diagram, project.Fs); err != nil {
//				panic(err)
//			}
//		}
//	}
//}

func LoadMarkdownTemplates(markdowns ...string) ([]*template.Template, error) {
	templates := make([]*template.Template, 0, len(markdowns))
	for _, markdownString := range markdowns {
		newTemplate, err := template.New("").Parse(markdownString)
		if err != nil {
			return nil, err
		}
		templates = append(templates, newTemplate)
	}
	return templates, nil
}

//func GenerateMarkdown(Index map[string]*IndexMarkdown, output string, fs afero.Fs) {
//	IndexTemplate, err := template.New("markdown").Parse(IndexMarkdownTemplate)
//	if err != nil {
//		panic(err)
//	}
//	AppTemplate, err := template.New("markdown").Parse(AppMarkdownTemplate)
//	if err != nil {
//		panic(err)
//	}
//	EmbededSvgTemplate, err := template.New("markdown").Parse(embededSvgTemplate)
//	if err != nil {
//		panic(err)
//	}
//	README, err := fs.Create(output + "/README.md")
//	fmt.Println("Creating", README.Name())
//	if err != nil {
//		panic(err)
//	}
//	if err := IndexTemplate.Execute(README, Index); err != nil {
//		panic(err)
//	}
//
//	for _, key := range alphabeticalIndex(Index) {
//		row := Index[key]
//		README, err := fs.Create(path.Join(output, row.PackageRelLink))
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("Creating", README.Name())
//		if err := AppTemplate.Execute(README, row.EndPoints); err != nil {
//			panic(err)
//		}
//		for _, Endpoint := range row.EndPoints {
//			embededSvg, err := fs.Create(path.Join(output, Endpoint.Package, Endpoint.EndpointName+ext+".md"))
//			if err != nil {
//				panic(err)
//			}
//			if err := EmbededSvgTemplate.Execute(embededSvg, Endpoint); err != nil {
//				panic(err)
//			}
//		}
//	}
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

func alphabeticalIndex(m map[string]*IndexMarkdown) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalApps(m map[string]*sysl.Application) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalEndpoints(m map[string]*sysl.Endpoint) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func alphabeticalPackage(m map[string]*Package) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
