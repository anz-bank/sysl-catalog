package templategeneration

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"

	"github.com/anz-bank/sysl/pkg/sequencediagram"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"

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
	SequenceDiagrams    []*SequenceDiagram // map[appName + endpointName]
	DataModelDiagrams   []*Diagram
}

func (p Package) RegisterIntegrationDiagrams(m *sysl.Module) {

}

func (p Package) RegisterDataModelDiagrams(m *sysl.Module) {

}

// AlphabeticalRows returns an alphabetically sorted list of packages of any project.
func (p Project) AlphabeticalRows() []*Package {
	packages := make([]*Package, 0, len(p.Packages))
	for _, key := range AlphabeticalPackage(p.Packages) {
		packages = append(packages, p.Packages[key])
	}
	return packages
}

// RegisterSequenceDiagrams creates sequence Diagrams from the sysl Module in Project.
func (p Project) RegisterSequenceDiagrams() error {
	for _, key := range AlphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		packageName, appName := GetAppPackageName(app)
		if syslutil.HasPattern(app.Attrs, "ignore") {
			p.Log.Infof("Skipping application %s", app.Name)
			continue
		}
		for _, key2 := range AlphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			if syslutil.HasPattern(endpoint.Attrs, "ignore") {
				p.Log.Infof("Skipping application %s", app.Name)
				continue
			}
			packageD := p.Packages[packageName]
			diagram, err := packageD.SequenceDiagramFromEndpoint(appName, endpoint)
			if err != nil {
				return err
			}
			if description := app.GetAttrs()["description"]; description != nil {
				diagram.AppComment = description.GetS()
			}
			p.Packages[packageName].SequenceDiagrams = append(packageD.SequenceDiagrams, diagram)
			if p.Packages[packageName].DataModelDiagrams == nil {
				p.Packages[packageName].DataModelDiagrams = []*Diagram{}
			}
		}
	}
	return nil
}

func (p Project) GenerateEndpointDataModel(parentAppName string, t *sysl.Type) string {
	pl := &datamodelCmd{}
	pl.Project = ""
	p.Fs.MkdirAll(pl.Output, os.ModePerm)
	pl.Direct = true
	pl.ClassFormat = "%(classname)"
	spclass := sequencediagram.ConstructFormatParser("", pl.ClassFormat)
	var stringBuilder strings.Builder
	dataParam := &catalogdiagrams.DataModelParam{}
	dataParam.Mod = p.Module

	v := datamodeldiagram.MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
	vNew := &catalogdiagrams.DataModelView{
		DataModelView: *v,
	}
	return vNew.GenerateDataView(dataParam, parentAppName, t, p.Module)
}
func sanitiseOutputName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "/", "")
}

// SequenceDiagramFromEndpoint generates a sequence diagram from a sysl endpoint
func (p Package) SequenceDiagramFromEndpoint(appName string, endpoint *sysl.Endpoint) (*SequenceDiagram, error) {
	call := fmt.Sprintf("%s <- %s", appName, endpoint.Name)
	var re = regexp.MustCompile(`(?m)(?:<:)(?:\s*\S+)`)
	var typeName string
	seq, err := CreateSequenceDiagram(p.Parent.Module, call)
	if err != nil {
		return nil, err
	}
	diagram := &SequenceDiagram{}
	diagram.Parent = &p
	diagram.AppName = appName
	diagram.EndpointName = endpoint.Name
	diagram.OutputFileName__ = sanitiseOutputName(appName + endpoint.Name)
	diagram.OutputDir = path.Join(p.Parent.Output, p.PackageName)
	diagram.DiagramString = seq
	diagram.Diagramtype = diagram_sequence
	diagram.OutputMarkdownFileName = p.Parent.OutputFileName
	diagram.OutputDataModel = []*Diagram{}
	diagram.InputDataModel = []*Diagram{}
	if description := endpoint.GetAttrs()["description"]; description != nil {
		diagram.EndpointComment = description.GetS()
	}
	for i, param := range endpoint.Param {
		if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
			if path := param.Type.GetTypeRef().GetRef().GetPath(); path != nil {
				appName = paramNameParts[0]
				typeName = path[0]
			} else {
				typeName = paramNameParts[0]
			}
		} else {
			typeName = paramNameParts[0]
		}
		typeref := &sysl.Type{
			Type: &sysl.Type_TypeRef{
				TypeRef: &sysl.ScopedRef{
					Ref: &sysl.Scope{Appname: &sysl.AppName{
						Part: []string{appName},
					},
						Path: []string{appName, typeName},
					},
				},
			},
		}
		newDiagram := &Diagram{
			Parent:           &p,
			OutputDir:        path.Join(p.Parent.Output, p.PackageName),
			AppName:          appName,
			DiagramString:    p.Parent.GenerateEndpointDataModel(appName, typeref),
			OutputFileName__: sanitiseOutputName(appName + endpoint.Name + "data-model-parameter" + strconv.Itoa(i)),
			EndpointName:     endpoint.Name,
		}
		diagram.InputDataModel = append(diagram.InputDataModel, newDiagram)
	}
	for i, stmnt := range endpoint.Stmt {
		if ret := stmnt.GetRet(); ret != nil {
			t := strings.ReplaceAll(re.FindString(ret.Payload), "<: ", "")
			if split := strings.Split(t, "."); len(split) > 1 {
				appName = split[0]
				typeName = split[1]
			} else {
				typeName = split[0]
			}
			typeref := &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Appname: &sysl.AppName{
							Part: []string{appName},
						},
							Path: []string{appName, typeName},
						},
					},
				},
			}
			newDiagram := &Diagram{
				Parent:           &p,
				OutputDir:        path.Join(p.Parent.Output, p.PackageName),
				AppName:          appName,
				DiagramString:    p.Parent.GenerateEndpointDataModel(appName, typeref),
				OutputFileName__: sanitiseOutputName(appName + endpoint.Name + "data-model-response" + strconv.Itoa(i)),
				EndpointName:     endpoint.Name,
			}
			diagram.OutputDataModel = append(diagram.OutputDataModel, newDiagram)
		}
	}
	return diagram, nil
}
