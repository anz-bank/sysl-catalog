// diagram-creation.go: all the methods attached to the generator object to be used in templating
package catalog

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/russross/blackfriday/v2"

	"github.com/anz-bank/protoc-gen-sysl/syslpopulate"

	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/sysl"
)

var (
	ofTypeSymbol = regexp.MustCompile(`(?m)(?:<:)(?:.*)`)
)

const (
	plantuml = iota
	mermaidjs
)

// CreateMarkdown is a wrapper function that also converts output markdown to html if in server mode
func (p *Generator) CreateMarkdown(t *template.Template, outputFileName string, i interface{}) error {
	var buf bytes.Buffer
	if err := t.Execute(&buf, i); err != nil {
		return err
	}
	if err := p.Fs.MkdirAll(path.Dir(outputFileName), os.ModePerm); err != nil {
		return err
	}
	f2, err := p.Fs.Create(outputFileName)
	if err != nil {
		return err
	}
	out := buf.Bytes()
	if p.Format == "html" && !p.DisableCss {
		raw := string(blackfriday.Run(out))
		raw = strings.ReplaceAll(raw, "README.md", p.OutputFileName)
		out = []byte(header + raw + style + endTags)
	}
	if _, err = f2.Write(out); err != nil {
		return err
	}
	return nil
}
func (p *Generator) CreateIntegrationDiagram(m *sysl.Module, title string, EPA bool) string {
	if p.Mermaid {
		return p.CreateIntegrationDiagramMermaid(m, title, EPA)
	}
	return p.CreateIntegrationDiagramPlantuml(m, title, EPA)
}

//TODO @ashwinsajiv: fill out this function to generate mermaid integration diagrams
func (p *Generator) CreateIntegrationDiagramMermaid(m *sysl.Module, title string, EPA bool) string {
	return "" //p.CreateFile("plantumlString", mermaidjs, "title", "integration.Output+p.Ext")
}

// CreateIntegrationDiagram creates an integration diagram and returns the filename
func (p *Generator) CreateIntegrationDiagramPlantuml(m *sysl.Module, title string, EPA bool) string {
	type intsCmd struct {
		diagrams.Plantumlmixin
		cmdutils.CmdContextParamIntgen
	}
	integration := intsCmd{}
	projectApp := createProjectApp(m.Apps)
	project := "__TEMP__"
	defer delete(p.Module.GetApps(), project)
	p.Module.GetApps()[project] = projectApp
	integration.Project = project
	integration.Output = "integration" + TernaryOperator(EPA, "EPA", "").(string)
	integration.Title = title

	integration.EPA = EPA
	integration.Clustered = true
	result, err := integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, p.Module, logrus.New())
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	plantumlString := result[integration.Output]
	return p.CreateFile(plantumlString, plantuml, title, integration.Output+p.Ext)
}

// CreateSequenceDiagram creates an sequence diagram and returns the filename
func (p *Generator) CreateSequenceDiagram(appName string, endpoint *sysl.Endpoint) string {
	m := p.Module
	call := fmt.Sprintf("%s <- %s", appName, endpoint.GetName())
	plantumlString, err := CreateSequenceDiagram(m, call)
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	packageName, _ := GetAppPackageName(p.Module.GetApps()[appName])
	return p.CreateFile(plantumlString, plantuml, packageName, appName, endpoint.GetName()+p.Ext)
}

// CreateParamDataModel creates a parameter data model and returns a filename
func (p *Generator) CreateParamDataModel(app *sysl.Application, param *sysl.Param) string {
	var appName, typeName string
	appName, typeName = GetAppTypeName(param)
	if appName == "" {
		appName = path.Join(app.Name.GetPart()...)
	}
	packageName, _ := GetAppPackageName(p.Module.GetApps()[appName])
	relatedTypes := catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: NewTypeRef(appName, typeName)}, p.Module)
	plantumlString := catalogdiagrams.GenerateDataModel(appName, relatedTypes)
	return p.CreateFile(plantumlString, plantuml, packageName, appName+p.Ext)
}

// GetReturnType converts an application and a param into a type, useful for getting attributes.
func (p *Generator) GetParamType(app *sysl.Application, param *sysl.Param) *sysl.Type {
	var appName, typeName string
	appName, typeName = GetAppTypeName(param)
	if appName == "" {
		appName = path.Join(app.GetName().GetPart()...)
	}
	return p.Module.Apps[appName].GetTypes()[typeName]
}

// GetReturnType converts an endpoint and a statement into a type, useful for getting attributes.
func (p *Generator) GetReturnType(endpoint *sysl.Endpoint, stmnt *sysl.Statement) *sysl.Type {
	var appName, typeName string
	ret := stmnt.GetRet()
	if ret == nil {
		return nil
	}
	t := strings.ReplaceAll(ofTypeSymbol.FindString(ret.GetPayload()), "<: ", "")
	if strings.Contains(t, "sequence of") {
		t = strings.ReplaceAll(t, "sequence of ", "")
	}
	if split := strings.Split(t, "."); len(split) > 1 {
		appName = split[0]
		typeName = split[1]
	} else {
		typeName = split[0]
	}
	if appName == "" {
		appName = strings.Join(endpoint.GetSource().GetPart(), "")
	}
	return p.Module.GetApps()[appName].GetTypes()[typeName]
}

// CreateReturnDataModel creates a return data model and returns a filename, or empty string if it wasn't a return statement.
func (p *Generator) CreateReturnDataModel(stmnt *sysl.Statement, endpoint *sysl.Endpoint) string {
	var sequence bool
	var typeref *sysl.Type
	var appName, typeName string
	ret := stmnt.GetRet()
	if ret == nil {
		return ""
	}
	t := strings.ReplaceAll(ofTypeSymbol.FindString(ret.Payload), "<: ", "")
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
		newSequenceName := endpoint.GetName() + "ReturnVal"
		newAppName := strings.Join(endpoint.GetSource().GetPart(), "")
		defer delete(p.Module.GetApps()[newAppName].GetTypes(), newSequenceName)
		p.Module.GetApps()[newAppName].GetTypes()[newSequenceName] = &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{
					AttrDefs: map[string]*sysl.Type{"sequence": {Type: &sysl.Type_Sequence{
						Sequence: syslpopulate.NewType(typeName, appName)},
					},
					},
				},
			},
		}
		typeref = NewTypeRef(appName, newSequenceName)
	} else {
		typeref = NewTypeRef(appName, typeName)
	}
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	return p.CreateTypeDiagram(p.Module.GetApps()[appName], typeName, typeref, true)
}

// CreateTypeDiagram creates a data model diagram and returns the filename
func (p *Generator) CreateTypeDiagram(app *sysl.Application, typeName string, t *sysl.Type, recursive bool) string {
	m := p.Module
	appName := strings.Join(app.Name.Part, "")
	typeref := NewTypeRef(appName, typeName)
	var plantumlString string
	if recursive {
		relatedTypes := catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: typeref}, m)
		plantumlString = catalogdiagrams.GenerateDataModel(appName, relatedTypes)
	} else {
		plantumlString = catalogdiagrams.GenerateDataModel(appName, map[string]*sysl.Type{typeName: t})
	}
	if _, ok := p.Module.GetApps()[appName]; !ok {
		return ""
	}
	packageName, _ := GetAppPackageName(p.Module.GetApps()[appName])
	return p.CreateFile(plantumlString, plantuml, packageName, appName, typeName+TernaryOperator(recursive, "", "simple").(string)+p.Ext)
}

// CreateFileName returns the absolute and relative filepaths
func CreateFileName(dir string, elems ...string) (string, string) {
	absolutefileName := path.Join(append([]string{dir}, Map(elems, SanitiseOutputName)...)...)
	relativefileName := strings.Replace(absolutefileName, dir+"/", "", 1)
	return absolutefileName, relativefileName
}

// CreateFile registers a file that needs to be created in p, or returns the embedded img tag if in server mode
func (p *Generator) CreateFile(contents string, diagramType int, absolute string, elems ...string) string {
	fileName, relativeFilepath := CreateFileName(path.Join(p.NextOutputDir, absolute), elems...)
	var fileContents string
	var targetMap map[string]string
	var err error
	switch diagramType {
	case plantuml:
		fileContents, err = PlantUMLURL(p.PlantumlService, contents)
		targetMap = p.FilesToCreate
	case mermaidjs:
		fileContents = contents
		targetMap = p.MermaidFilesToCreate
	default:
		panic("Wrong diagram type specified")
	}
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	// if p.ImageTags: return image tag from plantUML service
	if p.ImageTags && diagramType == plantuml {
		return fileContents
	}
	targetMap[fileName] = fileContents
	return relativeFilepath
}

// GenerateDataModel generates a data model for all of the types in app
func (p *Generator) GenerateDataModel(app *sysl.Application) string {
	appName := strings.Join(app.GetName().GetPart(), "")
	plantumlString := catalogdiagrams.GenerateDataModel(appName, app.GetTypes())
	if _, ok := p.Module.GetApps()[appName]; !ok {
		return ""
	}
	packageName, _ := GetAppPackageName(app)
	return p.CreateFile(plantumlString, plantuml, packageName, appName, "types"+p.Ext)
}

// CreateQueryParamDataModel returns a Query Parameter data model filename.
func (p *Generator) CreateQueryParamDataModel(CurrentAppName string, param *sysl.Endpoint_RestParams_QueryParam) string {
	var typeName, appName string
	var parsedType *sysl.Type
	switch param.GetType().GetType().(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.GetType()
		typeName = param.GetName()
	case *sysl.Type_TypeRef:
		appName, typeName = GetAppTypeName(param)
		if appName == "" {
			appName = CurrentAppName
		}
		parsedType = NewTypeRef(appName, typeName)
	}
	if _, ok := p.Module.GetApps()[appName]; !ok {
		return ""
	}
	return p.CreateTypeDiagram(p.Module.GetApps()[appName], typeName, parsedType, true)
}

// CreateQueryParamDataModel returns a Path Parameter data model filename.
func (p *Generator) CreatePathParamDataModel(CurrentAppName string, param *sysl.Endpoint_RestParams_QueryParam) string {
	var typeName, appName string
	var parsedType *sysl.Type
	switch param.Type.Type.(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.GetType()
		typeName = param.GetName()
	case *sysl.Type_TypeRef:
		appName, typeName = GetAppTypeName(param)
		if appName == "" {
			appName = CurrentAppName
		}
		parsedType = NewTypeRef(appName, typeName)
	}
	if _, ok := p.Module.GetApps()[appName]; !ok {
		return ""
	}
	return p.CreateTypeDiagram(p.Module.GetApps()[appName], typeName, parsedType, true)
}

func (p *Generator) getProjectApp(m *sysl.Module) (*sysl.Application, map[string]string) {
	includedProjects := Filter(
		SortedKeys(m.Apps),
		func(i string) bool {
			return syslutil.HasPattern(m.GetApps()[i].GetAttrs(), "project") &&
				m.GetApps()[i].SourceContext.File == p.SourceFileName
		},
	)
	if len(includedProjects) > 0 {
		set := make(map[string]string)
		app := m.GetApps()[includedProjects[0]]
		for _, e := range app.GetEndpoints() {
			if syslutil.HasPattern(e.GetAttrs(), "ignore") {
				continue
			}
			for _, e2 := range e.GetStmt() {
				set[e2.GetAction().GetAction()] = e.Name
			}
		}
		return m.GetApps()[includedProjects[0]], set
	}
	return nil, nil
}

// ModuleAsMacroPackage returns "macro packages" that map to the endpoints on the "project" application
/*
project[~project]: <-- first map
	FirstDivision: <-- first key
		package1 <--- second map["package1"]*sysl.Module
		package2 <--- second map["package2"]*sysl.Module
	SecondDivision:
		package3
*/

func (p *Generator) ModuleAsMacroPackage(m *sysl.Module) map[string]map[string]*sysl.Module {
	packages := make(map[string]map[string]*sysl.Module)
	_, includedProjects := p.getProjectApp(m)
	for _, key := range SortedKeys(m.GetApps()) {
		app := m.GetApps()[key]
		packageName := Attribute(app, "package")
		if packageName == "" {
			packageName = key
		}
		// what endpoint on the "project app" are we on?
		projectEndpoint, ok := includedProjects[packageName]
		if len(includedProjects) > 0 && !ok || (projectEndpoint == "") {
			continue
		}
		if syslutil.HasPattern(app.GetAttrs(), "ignore") || syslutil.HasPattern(app.GetAttrs(), "project") {
			continue
		}
		if _, ok := packages[projectEndpoint]; !ok {
			packages[projectEndpoint] = make(map[string]*sysl.Module)
		}
		if _, ok := packages[projectEndpoint][packageName]; !ok {
			packages[projectEndpoint][packageName] = &sysl.Module{Apps: map[string]*sysl.Application{}}
		}

		packages[projectEndpoint][packageName].GetApps()[strings.Join(app.GetName().GetPart(), "")] = app
	}
	return packages
}

func (p *Generator) ModuleAsPackages(m *sysl.Module) map[string]*sysl.Module {
	packages := make(map[string]*sysl.Module)
	_, includedProjects := p.getProjectApp(m)
	for _, key := range SortedKeys(m.GetApps()) {
		app := m.GetApps()[key]
		packageName := Attribute(app, "package")
		if packageName == "" {
			packageName = key
		}
		if _, ok := includedProjects[packageName]; len(includedProjects) > 0 && !ok {
			continue
		}
		if syslutil.HasPattern(app.GetAttrs(), "ignore") || syslutil.HasPattern(app.GetAttrs(), "project") {
			continue
		}
		if _, ok := packages[packageName]; !ok {
			packages[packageName] = &sysl.Module{Apps: map[string]*sysl.Application{}}
		}
		packages[packageName].GetApps()[strings.Join(app.GetName().GetPart(), "")] = app
	}
	return packages
}
