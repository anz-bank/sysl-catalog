package main

import (
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

// Project is the top level in the hierarchy of markdown generation
type Project struct {
	Title           string
	Output          string
	PlantumlService string
	OutputFileName  string
	Packages        map[string]*Package //Packages are the rows of the top level markdown
	Fs              afero.Fs
	Module          *sysl.Module
	ProjectTempl    *template.Template // Templ is used to template the Project struct
	PackageTempl    *template.Template // PackageTempl is passed down to all Packages
	EmbededTempl    *template.Template // This is passed down to all Diagrams
}

type Package struct {
	Parent              *Project
	OutputDir           string
	PackageName         string
	OutputFile          string
	IntegrationDiagrams []*Diagram
	SequenceDiagrams    []*Diagram
	DataModelDiagrams   []*Diagram
}

type Diagram struct {
	Parent                 *Package
	OutputDir              string
	OutputFileName         string
	Diagram                string
	OutputMarkdownFileName string
	Diagramtype            int
}

// NewProject generates a Project Markdwon object for all a sysl module
func NewProject(title, output, plantumlservice string, fs afero.Fs, module *sysl.Module) *Project {

	p := Project{
		Title:           title,
		Output:          output,
		Fs:              fs,
		Module:          module,
		Packages:        map[string]*Package{},
		PlantumlService: plantumlservice,
		OutputFileName:  "README.md",
	}
	p.initPackage()
	p.InsertSequenceDiagram()
	return &p
}

// RegisterTemplates registers templates for the project to use
func (p *Project) RegisterTemplates(projectTemplateString, packageTemplateString, embededTemplateString string) *Project {
	templates, err := LoadMarkdownTemplates(projectTemplateString, packageTemplateString, embededTemplateString)
	if err != nil {
		panic(err)
	}
	projectTemplate, packageTemplate, embededTemplate := templates[0], templates[1], templates[2]
	if err != nil {
		panic(err)
	}
	p.ProjectTempl = projectTemplate
	p.PackageTempl = packageTemplate
	p.EmbededTempl = embededTemplate
	return p
}
func (p *Project) initPackage() {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		packageName, _ := getPackageName(app)
		newPackage, ok := p.Packages[packageName]
		if !ok {
			newPackage = &Package{
				Parent:      p,
				PackageName: packageName,
				OutputDir:   path.Join(p.Output, packageName),
				OutputFile:  "README.md",
			}
		}
		p.Packages[packageName] = newPackage
	}
}

func getPackageName(app *sysl.Application) (string, string) {
	appName := strings.Join(app.Name.GetPart(), "")
	packageName := appName
	if attr := app.GetAttrs()["package"]; attr != nil {
		packageName = attr.GetS()
	}
	return packageName, appName
}

func (p Package) InsertIntegrationDiagrams(m *sysl.Module) {

}

func (p Package) InsertDataModelDiagrams(m *sysl.Module) {

}

func (p Project) AlphabeticalRows() []*Package {
	packages := make([]*Package, 0, len(p.Packages))
	for _, key := range alphabeticalPackage(p.Packages) {
		packages = append(packages, p.Packages[key])
	}
	return packages
}

func (p Project) InsertSequenceDiagram() {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		packageName, appName := getPackageName(app)
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		for _, key2 := range alphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			packageD := p.Packages[packageName]
			p.Packages[packageName].SequenceDiagrams = append(packageD.SequenceDiagrams, packageD.sequenceDiagramFromEndpoint(appName, endpoint))
		}
	}
}

func (p Package) sequenceDiagramFromEndpoint(appName string, endpoint *sysl.Endpoint) *Diagram {
	call := fmt.Sprintf("%s <- %s", appName, endpoint.Name)
	seq, err := CreateSequenceDiagram(p.Parent.Module, call)
	if err != nil {
		panic(err)
	}
	return &Diagram{
		Parent:                 &p,
		OutputFileName:         strings.ReplaceAll(call, " ", ""),
		OutputDir:              path.Join(p.Parent.Output, p.PackageName),
		Diagram:                seq,
		Diagramtype:            diagram_sequence,
		OutputMarkdownFileName: "README.md",
	}
}

const (
	diagram_integration = iota
	diagram_sequence
	diagram_datamodel
)
