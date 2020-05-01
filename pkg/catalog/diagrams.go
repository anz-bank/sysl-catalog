package catalog

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"
	"github.com/spf13/afero"

	"github.com/anz-bank/protoc-gen-sysl/syslpopulate"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/diagrams"
)

// PlantUMLDiagramString represents a plantuml diagram with other contextual info.
type Diagram struct {
	Parent                *Package
	Project               *Project
	Endpoint              *sysl.Endpoint
	App                   *sysl.Application
	Type                  *sysl.Type
	OutputDir             string
	PlantUMLDiagramString string
	OutputFileName__      string
	Diagramtype           string
}

func (d Diagram) Img() string {
	var fs afero.Fs
	var fastLoad bool
	if d.Parent != nil && d.Parent.Parent.Server {
		fs = d.Parent.Parent.Fs
		fastLoad = true
	} else if d.Project != nil && d.Project.Server {
		fs = d.Project.Fs
		fastLoad = true
	}
	if fastLoad {
		file, err := afero.ReadFile(fs, path.Join(d.OutputDir, d.OutputFileName__))
		if err != nil {
			panic(err)
		}
		return string(file)
	}

	return fmt.Sprintf(`<img src="%s" alt="image">`, d.OutputFileName__)
}

func (d Diagram) AppComment() string {
	if d.App == nil {
		return ""
	}
	if description := d.App.GetAttrs()["description"]; description != nil {
		return description.GetS()
	}
	return ""
}

func (d Diagram) EndpointComment() string {
	if d.Endpoint == nil {
		return ""
	}
	if description := d.Endpoint.GetAttrs()["description"]; description != nil {
		return description.GetS()
	}
	return ""
}

func (d Diagram) TypeComment() string {
	if d.Type == nil {
		return ""
	}
	if description := d.Type.GetAttrs()["description"]; description != nil {
		return description.GetS()
	}
	return ""
}

func (d Diagram) AppName() string {
	if d.App == nil {
		return ""
	}
	return strings.Join(d.App.GetName().GetPart(), ".")
}

func (d Diagram) EndpointName() string {
	if d.Endpoint == nil {
		return ""
	}
	return d.Endpoint.Name
}

func (d Diagram) EndpointNameWithoutSpaces() string {
	if d.Endpoint == nil {
		return ""
	}
	return strings.ReplaceAll(d.Endpoint.Name, " ", "")
}

// InputDataModel Generates request diagrams for the endpoint that's registered in s
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
			Parent:                s.Parent,
			Type:                  s.Parent.Parent.Module.Apps[appName].Types[typeName],
			OutputDir:             path.Join(s.Parent.Parent.Output, s.Parent.PackageName),
			App:                   s.Parent.Parent.Module.Apps[appName],
			PlantUMLDiagramString: s.Parent.Parent.GenerateDBDataModel(appName, catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{appName: typeref}, s.Parent.Parent.Module)),
			OutputFileName__:      sanitiseOutputName(appName+s.Endpoint.Name+"data-model-parameter"+strconv.Itoa(i)) + s.Parent.Parent.DiagramExt,
		}
		diagram = append(diagram, newDiagram)
	}
	return diagram
}

// OutputDataModel Generates return value diagrams for the endpoint that's registered in s
func (s Diagram) OutputDataModel() []*Diagram {
	appName := s.AppName()
	typeName := ""
	var diagram []*Diagram
	if s.Endpoint == nil {
		return nil
	}
	for i, stmnt := range s.Endpoint.Stmt {
		var sequence bool
		var typeref *sysl.Type
		if ret := stmnt.GetRet(); ret != nil {
			t := strings.ReplaceAll(re.FindString(ret.Payload), "<: ", "")
			if strings.Contains(t, "sequence of") {
				t = strings.ReplaceAll(t, "sequence of ", "")
				sequence = true
			}
			if split := strings.Split(t, "."); len(split) > 1 {
				appName = split[0]
				typeName = split[1]
			} else {
				typeName = split[0]
			}
			if sequence {
				s.App.Types[s.Endpoint.Name+"ReturnVal"] = &sysl.Type{

					Type: &sysl.Type_Tuple_{
						Tuple: &sysl.Type_Tuple{
							AttrDefs: map[string]*sysl.Type{"sequence": &sysl.Type{Type: &sysl.Type_Sequence{
								Sequence: syslpopulate.NewType(typeName, appName)},
							},
							},
						},
					},
				}

				typeref = &sysl.Type{
					Type: &sysl.Type_TypeRef{
						TypeRef: &sysl.ScopedRef{
							Ref: &sysl.Scope{Appname: &sysl.AppName{
								Part: []string{appName},
							},
								Path: []string{appName, s.Endpoint.Name + "ReturnVal"},
							},
						},
					},
				}
			} else {
				typeref = &sysl.Type{
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
			}

			newDiagram := &Diagram{
				Parent:                s.Parent,
				Type:                  s.Parent.Parent.Module.Apps[appName].Types[typeName],
				OutputDir:             path.Join(s.Parent.Parent.Output, s.Parent.PackageName),
				App:                   s.Parent.Parent.Module.Apps[appName],
				PlantUMLDiagramString: s.Parent.Parent.GenerateDBDataModel(appName, catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{appName: typeref}, s.Parent.Parent.Module)),
				OutputFileName__:      sanitiseOutputName(appName+s.Endpoint.Name+"data-model-response"+strconv.Itoa(i)) + s.Parent.Parent.DiagramExt,
			}
			diagram = append(diagram, newDiagram)
		}
	}
	return diagram
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func (sd *Diagram) GenerateDiagramAndMarkdown() error {
	var wg sync.WaitGroup
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__)
	wg.Add(1)
	go func() {
		diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.PlantUMLDiagramString, sd.Parent.Parent.Fs)
		wg.Done()
	}()
	for _, d := range sd.InputDataModel() {
		wg.Add(1)
		go func(s *Diagram) {
			outputFileName := path.Join(s.OutputDir, s.OutputFileName__)
			diagrams.OutputPlantuml(outputFileName, s.Parent.Parent.PlantumlService, s.PlantUMLDiagramString, s.Parent.Parent.Fs)
			wg.Done()
		}(d)

	}
	for _, d := range sd.OutputDataModel() {
		wg.Add(1)
		go func(s *Diagram) {
			outputFileName := path.Join(s.OutputDir, s.OutputFileName__)
			diagrams.OutputPlantuml(outputFileName, s.Parent.Parent.PlantumlService, s.PlantUMLDiagramString, s.Parent.Parent.Fs)
			wg.Done()
		}(d)
	}
	wg.Wait()
	return nil
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func GenerateDiagramAndMarkdown(sd *Diagram) error {
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__)
	return diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.PlantUMLDiagramString, sd.Parent.Parent.Fs)
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

type datamodelCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamDatagen
}

type intsCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamIntgen
}

func createProjectApp(Apps map[string]*sysl.Application) *sysl.Application {
	app := syslpopulate.NewApplication("")
	app.Endpoints = make(map[string]*sysl.Endpoint)
	app.Endpoints["_"] = syslpopulate.NewEndpoint("_")
	app.Endpoints["_"].Stmt = []*sysl.Statement{}
	for key, _ := range Apps {
		app.Endpoints["_"].Stmt = append(app.Endpoints["_"].Stmt, syslpopulate.NewStringStatement(key))
	}
	return app
}

func (p *Project) CreateIntegrationDiagrams(title, output string, projectApp *sysl.Application, EPA bool) (*Diagram, error) {
	m := *p.Module
	if projectApp.Attrs == nil {
		projectApp.Attrs = make(map[string]*sysl.Attribute)
	}
	if _, ok := projectApp.Attrs["appfmt"]; !ok {
		projectApp.Attrs["appfmt"] = &sysl.Attribute{
			Attribute: &sysl.Attribute_S{S: "%(appname)"},
		}
	}
	newMap := make(map[string]*sysl.Application)
	for k, v := range m.Apps {
		newMap[k] = v
	}
	m.Apps = newMap
	m.Apps[title] = projectApp

	intType := ""
	if EPA {
		intType = "EPA"
	}
	integration := intsCmd{}
	p.Fs.MkdirAll(output, os.ModePerm)
	integration.Output = path.Join(output, title+intType+"_integration") + p.DiagramExt
	integration.Title = title
	integration.Project = title
	integration.EPA = EPA
	integration.Clustered = true
	result, err := integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, &m, p.Log)
	if err != nil {
		return nil, err
	}
	if err := integration.GenerateFromMap(result, p.Fs); err != nil {
		return nil, err
	}

	bytes, err := afero.ReadFile(p.Fs, integration.Output)
	if err != nil {
		return nil, err
	}

	return &Diagram{
		Parent:                nil,
		Project:               p,
		OutputDir:             output,
		App:                   projectApp,
		PlantUMLDiagramString: string(bytes),
		OutputFileName__:      path.Join(title+intType+"_integration") + p.DiagramExt,
	}, nil
}
