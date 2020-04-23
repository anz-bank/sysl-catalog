package templategeneration

import (
	"fmt"
	"path"
	"regexp"

	"github.com/anz-bank/sysl/pkg/sysl"
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
