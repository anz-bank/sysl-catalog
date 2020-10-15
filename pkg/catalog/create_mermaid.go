package catalog

import (
	"github.com/anz-bank/sysl/pkg/mermaid/datamodeldiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/endpointanalysisdiagram"
	integration "github.com/anz-bank/sysl/pkg/mermaid/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
)

func (p *Generator) IntegrationMermaid(m *sysl.Module, title string, EPA bool) string {
	var result string
	var err error
	var apps []string
	for app := range m.Apps {
		apps = append(apps, app)
	}
	if len(apps) == 0 {
		p.Log.Error("Empty Apps")
		return ""
	}
	mod := p.RootModule
	if EPA {
		result, err = endpointanalysisdiagram.GenerateMultipleAppEndpointAnalysisDiagram(mod, apps)

	} else {
		result, err = integration.GenerateMultipleAppIntegrationDiagram(mod, apps)
	}
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	return result
}

func (p *Generator) SequenceMermaid(appName string, endpoint *sysl.Endpoint) string {
	m := p.RootModule
	result, error := sequencediagram.GenerateSequenceDiagram(m, appName, endpoint.GetName())
	if error != nil {
		p.Log.Error(error)
		return ""
	}
	return result
}

func (p *Generator) DataModelReturnMermaid(appName string, stmnt *sysl.Statement, endpoint *sysl.Endpoint) string {
	appName, typeName, _, _ := p.ExtractReturnInfo(appName, stmnt, endpoint)
	return p.DataModelMermaid(appName, typeName)
}

func (p *Generator) DataModelMermaid(appName, typeName string) string {
	defer func() {
		if err := recover(); err != nil {
			p.Log.Errorf("error creating param data model: %s", err)
		}
	}()
	var mermaidString string
	var err error

	if appName == "primitive" {
		mermaidString += datamodeldiagram.GeneratePrimitive(typeName)
	} else {
		mermaidString, err = datamodeldiagram.GenerateDataDiagramWithMapper(p.Mapper, appName, typeName)
		if err != nil {
			return ""
		}
	}
	return mermaidString
}

func (p *Generator) DataModelAliasMermaid(app *sysl.Application, param Param) string {
	defer func() {
		if err := recover(); err != nil {
			p.Log.Errorf("error creating param data model: %s", err)
		}
	}()
	info, typeName, _, _ := p.ExtractTypeInfo(app, param)
	return p.DataModelMermaid(info, typeName)
}

func (p *Generator) DataModelAppMermaid(app *sysl.Application) string {
	s, _ := datamodeldiagram.GenerateFullDataDiagram(&sysl.Module{Apps: map[string]*sysl.Application{"_": app}})
	return s
}
