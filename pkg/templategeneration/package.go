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
	ext = ".svg"
)

var re = regexp.MustCompile(`(?m)(?:<:)(?:\s*\S+)`)

// Package is the second level where apps and endpoints are specified.
type Package struct {
	Parent           *Project
	OutputDir        string
	PackageName      string
	OutputFile       string
	SequenceDiagrams map[string][]*Diagram // map[appName][]*Diagram
	DatabaseModel    map[string]*Diagram
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
		if _, ok := p.Packages[packageName]; !ok {
			p.Packages[packageName] = &Package{Parent: &p}
		}
		if p.Packages[packageName].SequenceDiagrams == nil {
			p.Packages[packageName].SequenceDiagrams = make(map[string][]*Diagram)
			p.Packages[packageName].SequenceDiagrams[appName] = make([]*Diagram, 0, 0)
		}
		if syslutil.HasPattern(app.Attrs, "db") {
			if p.Packages[packageName].DatabaseModel == nil {
				p.Packages[packageName].DatabaseModel = make(map[string]*Diagram)
			}
			p.Packages[packageName].DatabaseModel[appName] = &Diagram{
				Parent:           p.Packages[packageName],
				App:              app,
				DiagramString:    p.GenerateDBDataModel(appName),
				OutputDir:        path.Join(p.Output, packageName),
				OutputFileName__: sanitiseOutputName(appName + "db"),
			}
		}
		if len(app.Endpoints) == 0 {
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

			p.Packages[packageName].SequenceDiagrams[appName] = append(packageD.SequenceDiagrams[appName], diagram)
		}

	}
	return nil
}

func (p Project) GenerateDBDataModel(parentAppName string) string {
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
		TypeMap:       p.Module.Apps[parentAppName].Types,
	}
	return vNew.GenerateDataView(dataParam, parentAppName, nil, p.Module)
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
		TypeMap:       map[string]*sysl.Type{},
	}
	return vNew.GenerateDataView(dataParam, parentAppName, t, p.Module)
}

func sanitiseOutputName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "/", "")
}

func (s Diagram) InputDataModel() []*Diagram {
	appName := s.AppName()
	typeName := ""
	var diagram []*Diagram
	if s.Endpoint == nil {
		return nil
	}
	for i, param := range s.Endpoint.Param {
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
			Parent:           s.Parent,
			OutputDir:        path.Join(s.Parent.Parent.Output, s.Parent.PackageName),
			App:              s.Parent.Parent.Module.Apps[appName],
			DiagramString:    s.Parent.Parent.GenerateEndpointDataModel(appName, typeref),
			OutputFileName__: sanitiseOutputName(appName + s.Endpoint.Name + "data-model-parameter" + strconv.Itoa(i)),
		}
		diagram = append(diagram, newDiagram)
	}
	return diagram
}

func (s Diagram) OutputDataModel() []*Diagram {
	appName := s.AppName()
	typeName := ""
	var diagram []*Diagram
	if s.Endpoint == nil {
		return nil
	}
	for i, stmnt := range s.Endpoint.Stmt {
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
				Parent:           s.Parent,
				OutputDir:        path.Join(s.Parent.Parent.Output, s.Parent.PackageName),
				App:              s.Parent.Parent.Module.Apps[appName],
				DiagramString:    s.Parent.Parent.GenerateEndpointDataModel(appName, typeref),
				OutputFileName__: sanitiseOutputName(appName + s.Endpoint.Name + "data-model-response" + strconv.Itoa(i)),
			}
			diagram = append(diagram, newDiagram)
		}
	}
	return diagram
}

// SequenceDiagramFromEndpoint generates a sequence diagram from a sysl endpoint
func (p Package) SequenceDiagramFromEndpoint(appName string, endpoint *sysl.Endpoint) (*Diagram, error) {
	call := fmt.Sprintf("%s <- %s", appName, endpoint.Name)
	seq, err := CreateSequenceDiagram(p.Parent.Module, call)
	if err != nil {
		return nil, err
	}
	diagram := &Diagram{
		Parent:                 &p,
		Endpoint:               endpoint,
		App:                    p.Parent.Module.Apps[appName],
		OutputDir:              path.Join(p.Parent.Output, p.PackageName),
		DiagramString:          seq,
		OutputFileName__:       sanitiseOutputName(appName + endpoint.Name),
		OutputMarkdownFileName: p.Parent.OutputFileName,
		Diagramtype:            diagram_sequence,
	}
	return diagram, nil
}
