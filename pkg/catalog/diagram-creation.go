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
	re = regexp.MustCompile(`(?m)(?:<:)(?:.*)`)
)

// CreateMarkdown is a wrapper function that also converts output markdown to html if in server mode
func (p *Generator) CreateMarkdown(t *template.Template, outputFileName string, i interface{}) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, i); err != nil {
		p.Log.Error(err)
		return
	}
	if err := p.Fs.MkdirAll(path.Dir(outputFileName), os.ModePerm); err != nil {
		p.Log.Error(err)
		return
	}
	f2, err := p.Fs.Create(outputFileName)
	if err != nil {
		p.Log.Error(err)
		return
	}
	out := buf.Bytes()
	if p.Format == "html" && !p.DisableCss {
		out = []byte(header + strings.ReplaceAll(string(blackfriday.Run(out)), "README.md", p.OutputFileName) + style + endTags)

	}
	if _, err = f2.Write(out); err != nil {
		p.Log.Error(err)
	}
}

// CreateIntegrationDiagram creates an integration diagram and returns the filename
func (p *Generator) CreateIntegrationDiagram(m *sysl.Module, title string, EPA bool) string {
	type intsCmd struct {
		diagrams.Plantumlmixin
		cmdutils.CmdContextParamIntgen
	}
	projectApp := createProjectApp(m.Apps)
	p.Module.Apps["__TEMP__"] = projectApp
	integration := intsCmd{}
	integration.Output = "integration" + TernaryOperator(EPA, "EPA", "").(string)
	integration.Title = title
	integration.Project = "__TEMP__"
	integration.EPA = EPA
	integration.Clustered = true
	result, err := integrationdiagram.GenerateIntegrations(&integration.CmdContextParamIntgen, p.Module, logrus.New())
	delete(p.Module.Apps, "__TEMP__")
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	return p.CreateFileName(result[integration.Output], title, integration.Output+".svg")
}

// CreateSequenceDiagram creates an sequence diagram and returns the filename
func (p *Generator) CreateSequenceDiagram(appName string, endpoint *sysl.Endpoint) string {
	m := p.Module
	call := fmt.Sprintf("%s <- %s", appName, endpoint.Name)
	seq, err := CreateSequenceDiagram(m, call)
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	packageName, _ := GetAppPackageName(p.Module.Apps[appName])
	return p.CreateFileName(seq, packageName, appName, endpoint.Name+".svg")
}

// CreateParamDataModel creates a parameter data model and returns a filename
func (p *Generator) CreateParamDataModel(app *sysl.Application, param *sysl.Param) string {
	var appName, typeName string
	appName = path.Join(app.Name.GetPart()...)
	if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
		if typeNameParts := param.Type.GetTypeRef().GetRef().GetPath(); typeNameParts != nil {
			appName = paramNameParts[0]
			typeName = typeNameParts[0]
		} else {
			typeName = paramNameParts[0]
		}
	} else {
		typeName = paramNameParts[0]
	}
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	packageName, _ := GetAppPackageName(p.Module.Apps[appName])
	relatedTypes := catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: NewTypeRef(appName, typeName)}, p.Module)
	return p.CreateFileName(catalogdiagrams.GenerateDataModel(appName, relatedTypes), packageName, appName+".svg")
}

// GetReturnType converts an application and a param into a type, useful for getting attributes.
func (p *Generator) GetParamType(app *sysl.Application, param *sysl.Param) *sysl.Type {
	var appName, typeName string
	appName = path.Join(app.Name.GetPart()...)
	if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
		if typeNamePath := param.Type.GetTypeRef().GetRef().GetPath(); typeNamePath != nil {
			appName = paramNameParts[0]
			typeName = typeNamePath[0]
		} else {
			typeName = paramNameParts[0]
		}
	} else {
		typeName = paramNameParts[0]
	}
	return p.Module.Apps[appName].Types[typeName]
}

// GetReturnType converts an endpoint and a statement into a type, useful for getting attributes.
func (p *Generator) GetReturnType(endpoint *sysl.Endpoint, stmnt *sysl.Statement) *sysl.Type {
	var appName, typeName string
	ret := stmnt.GetRet()
	if ret == nil {
		return nil
	}
	t := strings.ReplaceAll(re.FindString(ret.Payload), "<: ", "")
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
		appName = strings.Join(endpoint.Source.Part, "")
	}
	return p.Module.Apps[appName].Types[typeName]
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
		newSequenceName := endpoint.Name + "ReturnVal"
		newAppName := strings.Join(endpoint.Source.Part, "")
		defer delete(p.Module.Apps[newAppName].Types, newSequenceName)
		p.Module.Apps[newAppName].Types[newSequenceName] = &sysl.Type{
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
	return p.CreateTypeDiagram(p.Module.Apps[appName], typeName, typeref, true)
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
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	packageName, _ := GetAppPackageName(p.Module.Apps[appName])
	return p.CreateFileName(plantumlString, packageName, appName, typeName+TernaryOperator(recursive, "", "simple").(string)+".svg")
}

// CreateFileName registers a file that needs to be created in p, or returns the embedded img tag if in server mode
func (p *Generator) CreateFileName(contents string, absolute string, elems ...string) string {
	// if fastload: return image tag from plantuml service
	fileName := path.Join(Map(append([]string{absolute}, elems...), SanitiseOutputName)...)
	plantumlURL, err := PlantUMLURL(p.PlantumlService, contents)
	if err != nil {
		p.Log.Error(err)
		return ""
	}
	if p.ImageTags {
		return plantumlURL
	}
	p.FilesToCreate[fileName] = plantumlURL
	return strings.Replace(fileName, absolute+"/", "", 1) // if not fast load: return image filename
}

// GenerateDataModel generates a data model for all of the types in app
func (p *Generator) GenerateDataModel(app *sysl.Application) string {
	appName := strings.Join(app.Name.Part, "")
	plantumlString := catalogdiagrams.GenerateDataModel(appName, app.Types)
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	packageName, _ := GetAppPackageName(app)
	return p.CreateFileName(plantumlString, packageName, appName, "types"+".svg")
}

// CreateQueryParamDataModel returns a Query Parameter data model filename.
func (p *Generator) CreateQueryParamDataModel(CurrentAppName string, param *sysl.Endpoint_RestParams_QueryParam) string {
	var typeName, appName string
	var parsedType *sysl.Type
	switch param.Type.Type.(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.Type
		typeName = param.GetName()
	case *sysl.Type_TypeRef:
		if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
			if typeNameParts := param.Type.GetTypeRef().GetRef().GetPath(); typeNameParts != nil {
				appName = paramNameParts[0]
				typeName = typeNameParts[0]
			} else {
				typeName = paramNameParts[0]
			}
		} else {
			typeName = param.Type.GetTypeRef().GetRef().GetPath()[0]
		}
		if appName == "" {
			appName = CurrentAppName
		}
		parsedType = NewTypeRef(appName, typeName)
	}
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	return p.CreateTypeDiagram(p.Module.Apps[appName], typeName, parsedType, true)
}

// CreateQueryParamDataModel returns a Path Parameter data model filename.
func (p *Generator) CreatePathParamDataModel(CurrentAppName string, param *sysl.Endpoint_RestParams_QueryParam) string {
	var typeName, appName string
	var parsedType *sysl.Type
	switch param.Type.Type.(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.Type
		typeName = param.GetName()
	case *sysl.Type_TypeRef:
		if appNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(appNameParts) > 0 {
			if typeNameParts := param.Type.GetTypeRef().GetRef().GetPath(); typeNameParts != nil {
				appName = appNameParts[0]
				typeName = typeNameParts[0]
			} else {
				typeName = appNameParts[0]
			}
		} else {
			typeName = param.Type.GetTypeRef().GetRef().GetPath()[0]
		}
		if appName == "" {
			appName = CurrentAppName
		}
		parsedType = NewTypeRef(appName, typeName)
	}
	if _, ok := p.Module.Apps[appName]; !ok {
		return ""
	}
	return p.CreateTypeDiagram(p.Module.Apps[appName], typeName, parsedType, true)
}
