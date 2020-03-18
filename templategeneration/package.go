package templategeneration

import (
	"fmt"
	"path"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

const (
	md           = ".md"
	ext          = ".svg"
	pageFilename = "README.md"
)

// Package is the second level where apps and endpoints are specified.
type Package struct {
	Parent              *Project
	OutputDir           string
	PackageName         string
	OutputFile          string
	IntegrationDiagrams []*Diagram
	SequenceDiagrams    []*Diagram
	DataModelDiagrams   []*Diagram
}

func (p Package) RegisterIntegrationDiagrams(m *sysl.Module) {

}

func (p Package) RegisterDataModelDiagrams(m *sysl.Module) {

}

// AlphabeticalRows returns an alphabetically sorted list of packages of any project.
func (p Project) AlphabeticalRows() []*Package {
	packages := make([]*Package, 0, len(p.Packages))
	for _, key := range alphabeticalPackage(p.Packages) {
		packages = append(packages, p.Packages[key])
	}
	return packages
}

// RegisterSequenceDiagrams creates sequence Diagrams from the sysl Module in Project.
func (p Project) RegisterSequenceDiagrams() error {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		packageName, appName := GetAppPackageName(app)
		if syslutil.HasPattern(app.Attrs, "ignore") {
			p.Log.Infof("Skipping application %s", app.Name)
			continue
		}
		for _, key2 := range alphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			packageD := p.Packages[packageName]
			diagram, err := packageD.SequenceDiagramFromEndpoint(appName, endpoint)
			if err != nil {
				return err
			}
			p.Packages[packageName].SequenceDiagrams = append(packageD.SequenceDiagrams, diagram)
		}
	}
	return nil
}

// SequenceDiagramFromEndpoint generates a sequence diagram from a sysl endpoint
func (p Package) SequenceDiagramFromEndpoint(appName string, endpoint *sysl.Endpoint) (*Diagram, error) {
	call := fmt.Sprintf("%s <- %s", appName, endpoint.Name)
	seq, err := CreateSequenceDiagram(p.Parent.Module, call)
	if err != nil {
		return nil, err
	}
	return &Diagram{
		Parent:                 &p,
		AppName:                appName,
		EndpointName:           endpoint.Name,
		OutputFileName__:       appName + endpoint.Name,
		OutputDir:              path.Join(p.Parent.Output, p.PackageName),
		Diagram:                seq,
		Diagramtype:            diagram_sequence,
		OutputMarkdownFileName: pageFilename,
	}, nil
}
