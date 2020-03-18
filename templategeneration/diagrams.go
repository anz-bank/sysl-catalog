package templategeneration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/pkg/errors"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/diagrams"
)

// Diagram represents a plantuml diagram with other contextual info.
type Diagram struct {
	Parent                 *Package
	OutputDir              string
	AppName                string
	EndpointName           string
	Diagram                string
	OutputFileName__       string
	OutputMarkdownFileName string
	Diagramtype            string
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func GenerateDiagramAndMarkdown(sd *Diagram) error {
	fmt.Println(sd.OutputFileName__)
	if err := GenerateMarkdown(sd.OutputDir, sd.OutputFileName__+md, sd, sd.Parent.Parent.EmbededTempl, sd.Parent.Parent.Fs); err != nil {
		return err
	}
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName__+ext)
	return diagrams.OutputPlantuml(outputFileName, sd.Parent.Parent.PlantumlService, sd.Diagram, sd.Parent.Parent.Fs)
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

//func CreateDataModelDiagram(m *sysl.Module) (string, error) {
//	dataParam := &datamodeldiagram.DataModelParam{
//		Mod:   m,
//		App:   m.Apps["Apple"],
//		Title: "Apple",
//	}
//	v := datamodeldiagram.MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
//	outmap[outputDir] = v.GenerateDataView(dataParam)
//
//	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, m, logrus.New())
//	if err != nil {
//		return err
//	}
//	return p.GenerateFromMap(outmap, args.Filesystem)
//}

type intsCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamIntgen
}

func (p *Project) CreateIntegrationDiagrams() error {
	if _, ok := p.Module.Apps[p.Title]; !ok {
		return fmt.Errorf(
			"There must be a app with the same name as the input file:" +
				"'foo.sysl' must have a project named 'foo'")
	}
	integration := intsCmd{}
	integration.Output = path.Join(p.Output, p.Title+ext)
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
	p.RootLevelIntegrationDiagram = &Diagram{
		Parent:                 nil,
		OutputDir:              p.Output,
		AppName:                p.Title,
		EndpointName:           "",
		Diagram:                "", // Leave this empty because the diagram is already created
		OutputFileName__:       integration.Output,
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
	//p.RootLevelIntegrationDiagram = &Diagram{
	//	Parent:                 nil,
	//	OutputDir:              p.Output,
	//	AppName:                p.Title,
	//	EndpointName:           "",
	//	Diagram:                string(out),
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
