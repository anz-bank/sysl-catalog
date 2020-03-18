package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/spf13/afero"
)

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
	p := NewProject(m.String(), output, plantumlService, fs, m).RegisterTemplates(IndexMarkdownTemplate, AppMarkdownTemplate, embededSvgTemplate)
	p.GenerateMarkdownAndDiagrams()
}

func (p *Project) GenerateMarkdownAndDiagrams() {
	if err := GenerateMarkdown(p.Output, p.OutputFileName, p, p.ProjectTempl, p.Fs); err != nil {
		panic(err)
	}
	for _, key := range alphabeticalPackage(p.Packages) {
		row := p.Packages[key]
		if err := GenerateMarkdown(row.OutputDir, row.OutputFile, row, row.Parent.PackageTempl, p.Fs); err != nil {
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
	fmt.Println(sd.OutputDir, sd.OutputFileName+".md")
	if err := GenerateMarkdown(sd.OutputDir, sd.OutputFileName+".md", sd, sd.Parent.Parent.EmbededTempl, sd.Parent.Parent.Fs); err != nil {
		return err
	}
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName+ext)
	return diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.Diagram, sd.Parent.Parent.Fs)
}

func GenerateMarkdown(outputdir, fileName string, object interface{}, t *template.Template, fs afero.Fs) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, object); err != nil {
		return err
	}
	fs.MkdirAll(outputdir, os.ModePerm)
	return afero.WriteFile(fs, filepath.Join(outputdir, fileName), buf.Bytes(), os.ModePerm)
}

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
