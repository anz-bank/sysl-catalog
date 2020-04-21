package templategeneration

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/diagrams"
)

// DiagramString represents a plantuml diagram with other contextual info.
type Diagram struct {
	Parent                 *Package
	Endpoint               *sysl.Endpoint
	App                    *sysl.Application
	Type                   *sysl.Type
	OutputDir              string
	DiagramString          string
	OutputFileName__       string
	OutputMarkdownFileName string
	Diagramtype            string
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

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func (sd *Diagram) GenerateDiagramAndMarkdown() error {
	var wg sync.WaitGroup
	fmt.Println(sd.OutputFileName__)
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__+ext)
	wg.Add(1)
	go func() {
		diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.DiagramString, sd.Parent.Parent.Fs)
		wg.Done()
	}()
	for _, d := range sd.InputDataModel() {
		wg.Add(1)
		go func(s *Diagram) {
			outputFileName := path.Join(s.OutputDir, s.OutputFileName__+ext)
			diagrams.OutputPlantuml(outputFileName, s.Parent.Parent.PlantumlService, s.DiagramString, s.Parent.Parent.Fs)
			wg.Done()
		}(d)

	}
	for _, d := range sd.OutputDataModel() {
		wg.Add(1)
		go func(s *Diagram) {
			outputFileName := path.Join(s.OutputDir, s.OutputFileName__+ext)
			diagrams.OutputPlantuml(outputFileName, s.Parent.Parent.PlantumlService, s.DiagramString, s.Parent.Parent.Fs)
			wg.Done()
		}(d)
	}
	wg.Wait()
	return nil
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func GenerateDiagramAndMarkdown(sd *Diagram) error {
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__+ext)
	return diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.DiagramString, sd.Parent.Parent.Fs)
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

func (p *Project) CreateIntegrationDiagrams() error {
	projectApp, ok := p.Module.Apps[p.Title]
	if !ok {
		return fmt.Errorf(
			"There must be a app with the same name as the input file:"+
				"'%s.sysl' must have a project named '%s'", p.Title, p.Title)
	} else if _, ok := projectApp.Attrs["appfmt"]; !ok {
		projectApp.Attrs["appfmt"] = &sysl.Attribute{
			Attribute: &sysl.Attribute_S{S: "%(appname)"},
		}
	}
	integration := intsCmd{}
	integration.Output = path.Join(p.Output, p.Title+"_integration_EPA"+ext)
	integration.Title = p.Title
	integration.Project = p.Title
	integration.EPA = true
	integration.Clustered = true
	result, err := integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, p.Module, p.Log)
	if err != nil {
		return err
	}
	if err := integration.GenerateFromMap(result, p.Fs); err != nil {
		return err
	}
	p.RootLevelIntegrationDiagramEPA = &Diagram{
		Parent:                 nil,
		OutputDir:              p.Output,
		App:                    projectApp,
		DiagramString:          "", // Leave this empty because the diagram is already created
		OutputFileName__:       p.Title + "_integration_EPA" + ext,
		OutputMarkdownFileName: "",
		Diagramtype:            "integration",
	}
	integration.EPA = false
	integration.Output = path.Join(p.Output, p.Title+"_integration"+ext)
	result, err = integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, p.Module, p.Log)
	if err != nil {
		return err
	}
	if err := integration.GenerateFromMap(result, p.Fs); err != nil {
		return err
	}
	p.RootLevelIntegrationDiagram = &Diagram{
		Parent:                 nil,
		OutputDir:              p.Output,
		App:                    projectApp,
		DiagramString:          "", // Leave this empty because the diagram is already created
		OutputFileName__:       p.Title + "_integration" + ext,
		OutputMarkdownFileName: "",
		Diagramtype:            "integration",
	}
	return nil
}
