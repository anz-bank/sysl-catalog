// diagram-creation.go: all the methods attached to the generator object to be used in templating
package catalog

import (
	"bytes"
	"fmt"
	"github.com/anz-bank/sysl/pkg/mermaid/datamodeldiagram"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/anz-bank/protoc-gen-sysl/newsysl"
	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/mermaid/endpointanalysisdiagram"

	integration "github.com/anz-bank/sysl/pkg/mermaid/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/mermaid/sequencediagram"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
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
	var converted bytes.Buffer
	if p.Format == "html" && !p.DisableCss {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
		if err := md.Convert(out, &converted); err != nil {
			return errors.Wrap(err, "Error converting markdown to html:")
		}
		raw := string(converted.Bytes())
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
		return p.CreateIntegrationDiagramMermaid(m, EPA)
	}
	return p.CreateIntegrationDiagramPlantuml(m, title, EPA)
}

func (p *Generator) CreateIntegrationDiagramMermaid(m *sysl.Module, EPA bool) string {
	var result string
	var err error
	var apps []string
	for app, _ := range m.Apps {
		apps = append(apps, app)
	}
	mod := p.RootModule
	if EPA {
		result, err = endpointanalysisdiagram.GenerateEndpointAnalysisDiagram(mod)
		if err != nil {
			p.Log.Error(err)
			return ""
		}
	} else {
		result, err = integration.GenerateMultipleAppIntegrationDiagram(mod, apps)
	}
	output := "integration" + TernaryOperator(EPA, "EPA", "").(string)
	return p.CreateFile(result, mermaidjs, output+p.Ext)
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
	return p.CreateFile(plantumlString, plantuml, integration.Output+p.Ext)
}

func (p *Generator) CreateSequenceDiagram(appName string, endpoint *sysl.Endpoint) string {
	if p.Mermaid {
		return p.CreateSequenceDiagramMermaid(appName, endpoint)
	}
	return p.CreateSequenceDiagramPlantuml(appName, endpoint)
}

func (p *Generator) CreateSequenceDiagramMermaid(appName string, endpoint *sysl.Endpoint) string {
	m := p.RootModule
	result, error := sequencediagram.GenerateSequenceDiagram(m, appName, endpoint.GetName())
	if error != nil {
		p.Log.Error(error)
		return ""
	}
	return p.CreateFile(result, mermaidjs, appName, endpoint.GetName()+p.Ext)
}

// CreateSequenceDiagram creates an sequence diagram and returns the filename
func (p *Generator) CreateSequenceDiagramPlantuml(appName string, endpoint *sysl.Endpoint) string {
	m := p.RootModule
	call := fmt.Sprintf("%s <- %s", appName, endpoint.GetName())
	plantumlString, err := CreateSequenceDiagram(m, call)
	if err != nil {
		p.Log.Error("Error creating sequence diagram:", err)
		os.Exit(1)
		return ""
	}
	return p.CreateFile(plantumlString, plantuml, appName, endpoint.GetName()+p.Ext)
}

// CreateParamDataModel creates a parameter data model and returns a filename
func (p *Generator) CreateParamDataModel(app *sysl.Application, param Typer) string {
	var appName, typeName string
	var getRecursive bool
	appName, typeName = GetAppTypeName(param)
	if appName == "" {
		appName = path.Join(app.GetName().GetPart()...)
	}
	if appName == "primitive" {
		getRecursive = false
	} else {
		getRecursive = true
	}
	return p.CreateTypeDiagram(appName, typeName, param.GetType(), getRecursive)
}

// GetReturnType converts an application and a param into a type, useful for getting attributes.
func (p *Generator) GetParamType(app *sysl.Application, param *sysl.Param) *sysl.Type {
	var appName, typeName string
	appName, typeName = GetAppTypeName(param)
	if appName == "" {
		appName = path.Join(app.GetName().GetPart()...)
	}
	return p.RootModule.Apps[appName].GetTypes()[typeName]
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
	return p.RootModule.GetApps()[appName].GetTypes()[typeName]
}

// CreateReturnDataModel creates a return data model and returns a filename, or empty string if it wasn't a return statement.
func (p *Generator) CreateReturnDataModel(appname string, stmnt *sysl.Statement, endpoint *sysl.Endpoint) string {
	var sequence, getRecursive bool
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
		appName = appname
		typeName = split[0]
	}
	if sequence {
		newSequenceName := endpoint.GetName() + "ReturnVal"
		newAppName := appname
		defer delete(p.RootModule.GetApps()[newAppName].GetTypes(), newSequenceName)
		p.RootModule.GetApps()[newAppName].GetTypes()[newSequenceName] = &sysl.Type{
			Type: &sysl.Type_Tuple_{
				Tuple: &sysl.Type_Tuple{
					AttrDefs: map[string]*sysl.Type{"sequence": {Type: &sysl.Type_Sequence{
						Sequence: newsysl.Type(typeName, appName)},
					},
					},
				},
			},
		}
		typeref = NewTypeRef(newAppName, newSequenceName)
	} else {
		typeref = NewTypeRef(appName, typeName)
	}
	if _, ok := p.RootModule.Apps[appName]; !ok {
		return ""
	}
	if syslwrapper.IsPrimitive(typeName) {
		getRecursive = false
		typeref = syslwrapper.MakePrimitive(typeName)
	} else {
		getRecursive = true
	}
	return p.CreateTypeDiagram(appName, typeName, typeref, getRecursive)
}

func (p *Generator) CreateTypeDiagram(appName string, typeName string, t *sysl.Type, recursive bool) string {
	if p.Mermaid {
		return p.CreateTypeDiagramMermaid(appName, typeName, recursive)
	}
	return p.CreateTypeDiagramPlantuml(appName, typeName, t, recursive)
}

func (p *Generator) CreateTypeDiagramMermaid(appName string, typeName string, recursive bool) string {
	m := p.RootModule
	var mermaidString string
	var err error
	if appName != "" && typeName != "" {
		if appName != "primitive" {
			mermaidString, err = datamodeldiagram.GenerateDataDiagramWithAppAndType(m, appName, typeName)
			if err!=nil {
				return ""
			}
		} else {
			mermaidString = "%% AUTOGENERATED CODE -- DO NOT EDIT!\n\n"
			mermaidString +=  fmt.Sprintf("classDiagram\n class %s{\n  var %s\n }", appName, appName)
		}
	}
	return p.CreateFile(mermaidString, mermaidjs, appName, typeName+TernaryOperator(recursive, "", "simple").(string)+p.Ext)
}

// CreateTypeDiagram creates a data model diagram and returns the filename
// It handles recursively getting the related types, or for primitives, just returns the
func (p *Generator) CreateTypeDiagramPlantuml(appName string, typeName string, t *sysl.Type, recursive bool) string {
	m := p.RootModule
	var plantumlString string
	if recursive {
		relatedTypes := catalogdiagrams.RecursivelyGetTypes(appName, map[string]*sysl.Type{typeName: NewTypeRef(appName, typeName)}, m)
		plantumlString = catalogdiagrams.GenerateDataModel(appName, relatedTypes)
		if _, ok := p.RootModule.GetApps()[appName]; !ok {
				return ""
			}
		// Handle Empty
		if len(relatedTypes) == 0 {
			return ""
		}
	} else {
		plantumlString = catalogdiagrams.GenerateDataModel(appName, map[string]*sysl.Type{typeName: t})
	}
	return p.CreateFile(plantumlString, plantuml, appName, typeName+TernaryOperator(recursive, "", "simple").(string)+p.Ext)
}

// CreateFileName returns the absolute and relative filepaths
func CreateFileName(dir string, elems ...string) (string, string) {
	var absolutefilePath, filename string
	for i, e := range elems {
		if i == len(elems)-1 {
			filename = strings.ToLower(SanitiseOutputName(e))
			absolutefilePath = path.Join(absolutefilePath, filename)
			break
		}
		absolutefilePath = path.Join(absolutefilePath, SanitiseOutputName(e))
	}
	if dir != "" {
		dir += "/"
	}
	return path.Join(dir, absolutefilePath), dir
}

// CreateRedoc registers a file that needs to be created when the input source context has extension .json or .yaml
func (p *Generator) CreateRedoc(sourceContext *sysl.SourceContext, appName string) string {
	if !IsOpenAPIFile(sourceContext) || !p.Redoc {
		return ""
	}
	absPath, _ := CreateFileName(p.CurrentDir, appName+".redoc.html")
	p.RedocFilesToCreate[absPath] = BuildSpecURL(sourceContext)
	link, _ := CreateFileName("", appName+".redoc.html")
	return link
}

// CreateFile registers a file that needs to be created in p, or returns the embedded img tag if in server mode
func (p *Generator) CreateFile(contents string, diagramType int, elems ...string) string {
	absFilePath, currentDir := CreateFileName(p.CurrentDir, elems...)
	var targetMap map[string]string
	var err error
	switch diagramType {
	case plantuml:
		if !strings.Contains(p.PlantumlService, ".jar") {
			contents, err = PlantUMLURL(p.PlantumlService, contents)
		}
		targetMap = p.FilesToCreate
	case mermaidjs:
		targetMap = p.MermaidFilesToCreate
	default:
		panic("Wrong diagram type specified")
	}
	if err != nil {
		p.Log.Error("Error creating file:", err)
		os.Exit(1)
		return ""
	}
	newFileName := absFilePath
	for i := 0; ; i++ {
		if diagram, ok := targetMap[newFileName]; !ok || diagram == contents {
			break
		}
		newFileName = strings.ReplaceAll(absFilePath, p.Ext, strconv.Itoa(i)+p.Ext)
	}
	absFilePath = newFileName
	// if p.ImageTags: return image tag from plantUML service
	if p.ImageTags && diagramType == plantuml && !strings.Contains(p.PlantumlService, ".jar") {
		return contents
	}
	targetMap[absFilePath] = contents
	return strings.Replace(absFilePath, currentDir, "", 1)
}

// GenerateDataModel generates a data model for all of the types in app
func (p *Generator) GenerateDataModel(app *sysl.Application) string {
	appName := strings.Join(app.GetName().GetPart(), "")
	plantumlString := catalogdiagrams.GenerateDataModel(appName, app.GetTypes())
	if _, ok := p.RootModule.GetApps()[appName]; !ok {
		return ""
	}
	return p.CreateFile(plantumlString, plantuml, appName, "types"+p.Ext)
}

func (p *Generator) getProjectApp(m *sysl.Module) (*sysl.Application, map[string]string) {
	includedProjects := Filter(
		SortedKeys(m.Apps),
		func(i string) bool {
			return syslutil.HasPattern(m.GetApps()[i].GetAttrs(), "project") &&
				m.GetApps()[i].SourceContext.File == path.Base(p.SourceFileName)
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
	FirstDivision: <-- ["FirstDivision"]:
		package1   <-- all apps in this package will be contained in this module
		package2   <-- same with this one
	SecondDivision: <-- ["SecondDivision"]:
		package3    <-- You get the point
*/

func (p *Generator) ModuleAsMacroPackage(m *sysl.Module) map[string]*sysl.Module {
	packages := make(map[string]*sysl.Module)
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
			packages[projectEndpoint] = newsysl.Module()
		}
		if _, ok := packages[projectEndpoint]; !ok {
			packages[projectEndpoint] = &sysl.Module{Apps: map[string]*sysl.Application{}}
		}
		packages[projectEndpoint].GetApps()[strings.Join(app.GetName().GetPart(), "")] = app
	}
	return packages
}

// MacroPackages executes the markdown for a MacroPackage and returns a slice of the rows
func (p *Generator) MacroPackages(module *sysl.Module) []string {
	defer p.resetTempVars()
	MacroPackages := p.ModuleAsMacroPackage(module)
	for macroPackageName, macroPackage := range MacroPackages {
		fileName := markdownName(p.OutputFileName, macroPackageName)
		macroPackageFileName := path.Join(p.OutputDir, macroPackageName, fileName)
		p.CurrentDir = macroPackageName
		p.TempDir = macroPackageName // this is for p.Packages()
		p.Title = macroPackageName
		p.Links = map[string]string{"Back": "../" + p.OutputFileName}
		p.Module = macroPackage
		err := p.CreateMarkdown(p.Templates[1], macroPackageFileName, p)
		if err != nil {
			p.Log.Error("Error generating project table:", err)
			os.Exit(1)
		}
	}
	return SortedKeys(MacroPackages)
}

func (p *Generator) resetTempVars() {
	p.CurrentDir = p.TempDir
	p.TempDir = ""
}

// Packages executes the markdown for a package and returns a slice of the rows
func (p *Generator) Packages(m *sysl.Module) []string {
	defer p.resetTempVars()
	MacroPackages := p.ModuleAsPackages(m)
	for packageName, pkg := range MacroPackages {
		p.CurrentDir = path.Join(p.TempDir, packageName)
		fileName := markdownName(p.OutputFileName, packageName)
		fullOutputName := path.Join(p.OutputDir, p.CurrentDir, fileName)
		if err := p.CreateMarkdown(p.Templates[len(p.Templates)-1], fullOutputName, pkg); err != nil {
			p.Log.Error("error in generating "+fullOutputName, err)
		}
	}
	return SortedKeys(MacroPackages)
}

// ModuleAsPackages returns a map of [packagename]*sysl.Module
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
