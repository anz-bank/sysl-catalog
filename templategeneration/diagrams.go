package templategeneration

import (
	"fmt"
	"path"

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
	OutputFileName         string
	Diagram                string
	OutputMarkdownFileName string
	Diagramtype            string
}

// GenerateDiagramAndMarkdown generates diagrams and markdown for sysl diagrams.
func GenerateDiagramAndMarkdown(sd *Diagram) error {
	fmt.Println(sd.OutputDir, sd.OutputFileName+md)
	if err := GenerateMarkdown(sd.OutputDir, sd.OutputFileName+md, sd, sd.Parent.Parent.EmbededTempl, sd.Parent.Parent.Fs); err != nil {
		return err
	}
	outputFileName := path.Join(sd.OutputDir, sd.OutputFileName+ext)
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
