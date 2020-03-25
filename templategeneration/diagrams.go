package templategeneration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/anz-bank/sysl/pkg/datamodeldiagram"

	"github.com/pkg/errors"

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
	OutputDir              string
	AppName                string
	EndpointName           string
	DiagramString          string
	OutputFileName__       string
	OutputMarkdownFileName string
	Diagramtype            string
}

type SequenceDiagram struct {
	Diagram
	Endpoint        *sysl.Endpoint
	InputDataModel  []*Diagram
	OutputDataModel []*Diagram
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func (sd *SequenceDiagram) GenerateDiagramAndMarkdown() error {
	fmt.Println(sd.OutputFileName__)
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__+ext)
	diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.DiagramString, sd.Parent.Parent.Fs)
	for _, d := range sd.InputDataModel {
		outputFileName := path.Join(d.OutputDir, d.OutputFileName__+ext)
		diagrams.OutputPlantuml(outputFileName, d.Parent.Parent.PlantumlService, d.DiagramString, d.Parent.Parent.Fs)
	}
	for _, d := range sd.OutputDataModel {
		outputFileName := path.Join(d.OutputDir, d.OutputFileName__+ext)
		diagrams.OutputPlantuml(outputFileName, d.Parent.Parent.PlantumlService, d.DiagramString, d.Parent.Parent.Fs)
	}
	return nil
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func GenerateDiagramAndMarkdown(sd *Diagram) error {
	fmt.Println(sd.OutputFileName__)
	if err := GenerateMarkdown(sd.OutputDir, sd.OutputFileName__+md, sd, sd.Parent.Parent.EmbededTempl, sd.Parent.Parent.Fs); err != nil {
		return err
	}
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

func (p *Project) CreateDataModelDiagram() (string, error) {
	for name, m := range p.PackageModules {
		skip := true
		for _, this := range m.Apps {
			if len(this.Types) > 0 {
				skip = false
			}
		}
		if skip {
			continue
		}
		pl := &datamodelCmd{}
		pl.Project = ""
		pl.Output = path.Join(p.Output, name)
		p.Fs.MkdirAll(pl.Output, os.ModePerm)
		pl.Output += "/" + name + "_datamodel.svg"
		pl.Direct = true
		pl.ClassFormat = "%(classname)"
		outmap, err := datamodeldiagram.GenerateDataModels(&pl.CmdContextParamDatagen, m, logrus.New())
		if err != nil {
			return "", err
		}
		err = pl.GenerateFromMap(outmap, p.Fs)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

type intsCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamIntgen
}

func (p *Project) CreateIntegrationDiagrams() error {
	if _, ok := p.Module.Apps[p.Title]; !ok {
		return fmt.Errorf(
			"There must be a app with the same name as the input file:"+
				"'%ss.sysl' must have a project named '%s'", p.Title, p.Title)
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
		AppName:                p.Title,
		EndpointName:           "",
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
		AppName:                p.Title,
		EndpointName:           "",
		DiagramString:          "", // Leave this empty because the diagram is already created
		OutputFileName__:       p.Title + "_integration" + ext,
		OutputMarkdownFileName: "",
		Diagramtype:            "integration",
	}
	return nil
	// TODO: defer actual creation of file till later with something like this:
	//res := ""
	//for _, diag := range result {
	//	res = diag
	//}
	//
	//encoded, err := diagrams.DeflateAndEncode([]byte(res))
	//plantuml := fmt.Sprintf("%s/%s/%s", p.PlantumlService, ext, encoded)
	//out, err := SendHTTPRequest(plantuml)
	//if err != nil {
	//	return err
	//}
	//p.RootLevelIntegrationDiagramEPA = &DiagramString{
	//	Parent:                 nil,
	//	OutputDir:              p.Output,
	//	AppName:                p.Title,
	//	EndpointName:           "",
	//	DiagramString:                string(out),
	//	OutputFileName__:       integration.Output,
	//	OutputMarkdownFileName: "",
	//	Diagramtype:            "integration",
	//}
}

func SendHTTPRequest(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, errors.Errorf("Unable to create http request to %s, Error:%s", url, err.Error())
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Unable to read body.")
	}
	return out, nil
}
