package catalog

import (
	"fmt"
	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"os"
)

// IntegrationPlantuml creates an integration diagram and returns the plantuml string
func (p *Generator) IntegrationPlantuml(m *sysl.Module, title string, EPA bool) string {
	type intsCmd struct {
		diagrams.Plantumlmixin
		cmdutils.CmdContextParamIntgen
	}
	integration := intsCmd{}
	projectApp := createProjectApp(m.Apps)
	project := "__TEMP__"
	defer delete(p.RootModule.GetApps(), project)
	p.RootModule.GetApps()[project] = projectApp
	integration.Project = project
	integration.Output = "integration" + TernaryOperator(EPA, "EPA", "").(string)
	integration.Title = title
	integration.EPA = EPA
	integration.Clustered = true
	result, err := integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, p.RootModule, p.Log)
	if err != nil {
		p.Log.Error("Error creating integration diagram:", err)
		os.Exit(1)
	}
	plantumlString := result[integration.Output]
	return PlantUMLURL(p.PlantumlService, plantumlString)
}

// SequencePlantuml creates an sequence diagram and returns a plantuml url
func (p *Generator) SequencePlantuml(appName string, endpoint *sysl.Endpoint) string {
	m := p.RootModule
	call := fmt.Sprintf("%s <- %s", appName, endpoint.GetName())
	plantumlString, err := CreateSequenceDiagram(m, call)
	if err != nil {
		p.Log.Error("Error creating sequence diagram:", err)
		os.Exit(1)
		return ""
	}
	return PlantUMLURL(p.PlantumlService, plantumlString)

}

// DataModelReturnPlantuml returns a return datamodel as a plantuml url
func (p *Generator) DataModelReturnPlantuml(appName string, stmnt *sysl.Statement, endpoint *sysl.Endpoint) string {
	appName, typeName, t, recursive := p.ExtractReturnInfo(appName, stmnt, endpoint)
	if t == nil {
		return ""
	}
	return p.DataModelAliasPlantuml(appName, typeName, typeName, t, recursive)
}

// DataModelParamPlantuml creates a parameter data model and returns a plantuml url
func (p *Generator) DataModelParamPlantuml(app *sysl.Application, param Param) string {
	defer func() {
		if err := recover(); err != nil {
			p.Log.Errorf("error creating param data model: %s", err)
		}
	}()
	appName, typeName, aliasTypeName, getRecursive := p.ExtractTypeInfo(app, param)
	return p.DataModelAliasPlantuml(appName, typeName, aliasTypeName, param.GetType(), getRecursive)
}

// DataModelAppPlantuml generates a data model for all of the types in app and returns a plantuml url
func (p *Generator) DataModelAppPlantuml(app *sysl.Application) string {
	appName := GetAppNameString(app)
	plantumlString := catalogdiagrams.GenerateDataModel(appName, catalogdiagrams.FromSyslTypeMap(appName, app.GetTypes()))
	if _, ok := p.RootModule.GetApps()[appName]; !ok {
		return ""
	}
	return PlantUMLURL(p.PlantumlService, plantumlString)
}

func (p *Generator) DataModelPlantuml(appName, typeName string, t *sysl.Type, recursive bool) string {
	return p.DataModelAliasPlantuml(appName, typeName, typeName, t, recursive)
}

// DataModelAliasPlantuml generates a plantuml url for a alias data model
func (p *Generator) DataModelAliasPlantuml(appName, typeName, typeAlias string, t *sysl.Type, recursive bool) string {
	m := p.RootModule
	var plantumlString string
	if recursive {
		relatedTypes := catalogdiagrams.RecursivelyGetTypes(
			appName,
			map[string]*catalogdiagrams.TypeData{
				typeAlias: catalogdiagrams.NewTypeData(typeAlias, NewTypeRef(appName, typeName)),
			}, m,
		)
		plantumlString = catalogdiagrams.GenerateDataModel(appName, relatedTypes)
		if _, ok := p.RootModule.GetApps()[appName]; !ok {
			p.Log.Warnf("no app named %s", appName)
			return ""
		}
		// Handle Empty
		if len(relatedTypes) == 0 {
			return ""
		}
	} else {
		plantumlString = catalogdiagrams.GenerateDataModel(
			appName,
			map[string]*catalogdiagrams.TypeData{typeAlias: catalogdiagrams.NewTypeData(typeAlias, t)},
		)
	}
	return PlantUMLURL(p.PlantumlService, plantumlString)
}
