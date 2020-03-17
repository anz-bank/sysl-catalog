package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

type Project struct {
	Title           string
	Output          string
	Fs              afero.Fs
	Module          *sysl.Module
	rows            map[string]*Package
	PlantumlService string
}

type Package struct {
	Parent              *Project
	PackageName         string
	IntegrationDiagrams []*Diagram
	SequenceDiagrams    []*Diagram
	DataModelDiagrams   []*Diagram
}

func NewProject(title, output, plantumlservice string, fs afero.Fs, module *sysl.Module) Project {
	p := Project{
		Title:           title,
		Output:          output,
		Fs:              fs,
		Module:          module,
		rows:            map[string]*Package{},
		PlantumlService: plantumlservice,
	}
	p.initPackage()
	p.InsertSequenceDiagram()
	return p
}

func (p *Project) initPackage() {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		packageName, _ := getPackageName(app)
		newPackage, ok := p.rows[packageName]
		if !ok {
			newPackage = &Package{Parent: p, PackageName: packageName}
		}
		p.rows[packageName] = newPackage
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

func (p Project) InsertSequenceDiagram() {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		packageName, appName := getPackageName(app)
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		for _, key2 := range alphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			packageD := p.rows[packageName]
			p.rows[packageName].SequenceDiagrams = append(packageD.SequenceDiagrams, packageD.sequenceDiagramFromEndpoint(appName, endpoint))
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
		Parent:          &p,
		Name:            strings.ReplaceAll(call, " ", ""),
		OutputDirectory: path.Join(p.Parent.Output, p.PackageName),
		Diagram:         seq,
		Diagramtype:     diagram_sequence,
	}
}

type Diagram struct {
	Parent          *Package
	OutputDirectory string
	Name            string
	Diagram         string
	Diagramtype     int
}

const (
	diagram_integration = iota
	diagram_sequence
	diagram_datamodel
)
