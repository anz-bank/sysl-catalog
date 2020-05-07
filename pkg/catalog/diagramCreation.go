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
	"github.com/anz-bank/sysl/pkg/sequencediagram"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/sysl"
)

var re = regexp.MustCompile(`(?m)(?:<:)(?:.*)`)

// CreateMarkdown is a wrapper function that also converts output markdown to html if in server mode
func (p *Generator) CreateMarkdown(t *template.Template, outputFileName string, i interface{}) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, i); err != nil {
		panic(err)
	}
	if err := p.Fs.MkdirAll(path.Dir(outputFileName), os.ModePerm); err != nil {
		panic(err)
	}
	f2, err := p.Fs.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	out := buf.Bytes()
	if p.Format == "html" {
		out = []byte(header + string(blackfriday.Run(out)) + style + endTags)
	}
	f2.Write(out)

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
		panic(err)
	}
	packageName, _ := GetAppPackageName(p.Module.Apps[appName])
	return p.CreateFileName(seq, packageName, appName, endpoint.Name+".svg")
}

// CreateSequenceDiagram creates an sequence diagram and returns the sequence diagram string and any errors
func CreateSequenceDiagram(m *sysl.Module, call string) (string, error) {
	l := &cmdutils.Labeler{}
	p := &sequencediagram.SequenceDiagParam{}
	p.Endpoints = []string{call}
	p.AppLabeler = l
	p.EndpointLabeler = l
	p.Title = call
	return sequencediagram.GenerateSequenceDiag(m, p, logrus.New())
}

// CreateParamDataModel creates a parameter data model and returns a filename
func (p *Generator) CreateParamDataModel(app *sysl.Application, param *sysl.Param) string {
	var appName, typeName string
	appName = path.Join(app.Name.GetPart()...)
	if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
		if path := param.Type.GetTypeRef().GetRef().GetPath(); path != nil {
			appName = paramNameParts[0]
			typeName = path[0]
		} else {
			typeName = paramNameParts[0]
		}
	} else {
		typeName = paramNameParts[0]
	}

	typeref := NewTypeRef(appName, typeName)
	if _, ok := p.Module.Apps[appName]; ok {
		packageName, _ := GetAppPackageName(p.Module.Apps[appName])
		relatedTypes := catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: typeref}, p.Module)
		return p.CreateFileName(catalogdiagrams.GenerateDataModel(appName, relatedTypes), packageName, appName+".svg")
	}
	return ""
}

func (p *Generator) GetParamType(app *sysl.Application, param *sysl.Param) *sysl.Type {
	var appName, typeName string
	appName = path.Join(app.Name.GetPart()...)
	if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
		if path := param.Type.GetTypeRef().GetRef().GetPath(); path != nil {
			appName = paramNameParts[0]
			typeName = path[0]
		} else {
			typeName = paramNameParts[0]
		}
	} else {
		typeName = paramNameParts[0]
	}
	return p.Module.Apps[appName].Types[typeName]
}

// CreateReturnDataModel creates a return data model and returns a filename, or empty string if it wasn't a return statement
func (p *Generator) CreateReturnDataModel(stmnt *sysl.Statement, endpoint *sysl.Endpoint) string {
	var sequence bool
	var typeref *sysl.Type
	var appName, typeName string
	if ret := stmnt.GetRet(); ret != nil {
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
			p.Module.Apps[strings.Join(endpoint.Source.Part, "")].Types[endpoint.Name+"ReturnVal"] = &sysl.Type{
				Type: &sysl.Type_Tuple_{
					Tuple: &sysl.Type_Tuple{
						AttrDefs: map[string]*sysl.Type{"sequence": {Type: &sysl.Type_Sequence{
							Sequence: syslpopulate.NewType(typeName, appName)},
						},
						},
					},
				},
			}
			typeref = NewTypeRef(appName, endpoint.Name+"ReturnVal")
		} else {
			typeref = NewTypeRef(appName, typeName)
		}
		relatedReturnTypes := catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: typeref}, p.Module)
		// Don't generate diagrams for empty types e.g Types defined with ...
		if len(relatedReturnTypes) == 1 && relatedReturnTypes[appName+"."+typeName].Type == nil {
			return ""
		}
		if _, ok := p.Module.Apps[appName]; ok {
			packageName, _ := GetAppPackageName(p.Module.Apps[appName])
			return p.CreateFileName(catalogdiagrams.GenerateDataModel(appName, relatedReturnTypes), packageName, appName, typeName+".svg")
		}
	}
	return ""
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
	if _, ok := p.Module.Apps[appName]; ok {
		packageName, _ := GetAppPackageName(p.Module.Apps[appName])
		return p.CreateFileName(plantumlString, packageName, appName, typeName+TernaryOperator(recursive, "full", "simple").(string)+".svg")
	}
	return ""
}

// CreateFileName registers a file that needs to be created in p, or returns the embedded img tag if in server mode
func (p *Generator) CreateFileName(contents string, absolute string, elems ...string) string {
	// if fastload: return image from plantuml service
	fileName := path.Join(Map(append([]string{absolute}, elems...), SanitiseOutputName)...)
	if p.Server {
		encoded, err := diagrams.DeflateAndEncode([]byte(contents))
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s/%s/%s", p.PlantumlService, "svg", encoded)
	} else {
		p.FilesToCreate[fileName] = contents
	}
	// if not fast load: return image filename
	return strings.Replace(fileName, absolute+"/", "", 1)
}

// GenerateDataModel generates a data model for all of the types in app
func (p *Generator) GenerateDataModel(app *sysl.Application) string {
	appName := strings.Join(app.Name.Part, "")
	plantumlString := catalogdiagrams.GenerateDataModel(appName, app.Types)
	if _, ok := p.Module.Apps[appName]; ok {
		packageName, _ := GetAppPackageName(app)
		return p.CreateFileName(plantumlString, packageName, appName, "types"+".svg")
	}
	return ""
}

func (p *Generator) CreateQueryParamDataModel(param *sysl.Param) string {
	var typeName, appName string
	relatedTypes := make(map[string]*sysl.Type)
	var parsedType *sysl.Type
	switch param.Type.Type.(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.Type
		typeName = param.GetName()
		relatedTypes = map[string]*sysl.Type{appName + ":" + typeName: parsedType}
	case *sysl.Type_TypeRef:
		if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
			if path := param.Type.GetTypeRef().GetRef().GetPath(); path != nil {
				appName = paramNameParts[0]
				typeName = path[0]
			} else {
				typeName = paramNameParts[0]
			}
		} else {
			typeName = param.Type.GetTypeRef().GetRef().GetPath()[0]
		}
		parsedType = NewTypeRef(appName, typeName)
		relatedTypes = catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: parsedType}, p.Module)
	}
	if _, ok := p.Module.Apps[appName]; ok {
		packageName, _ := GetAppPackageName(p.Module.Apps[appName])
		return p.CreateFileName(catalogdiagrams.GenerateDataModel(appName, relatedTypes), packageName, appName, typeName+"full.svg")
	}
	return ""
}

func (p *Generator) CreatePathParamDataModel(param *sysl.Param) string {
	var typeName, appName string
	relatedTypes := make(map[string]*sysl.Type)
	var parsedType *sysl.Type
	switch param.Type.Type.(type) {
	case *sysl.Type_Primitive_:
		parsedType = param.Type
		typeName = param.GetName()
		relatedTypes = map[string]*sysl.Type{appName + ":" + typeName: parsedType}
	case *sysl.Type_TypeRef:
		if paramNameParts := param.Type.GetTypeRef().GetRef().GetAppname().GetPart(); len(paramNameParts) > 0 {
			if path := param.Type.GetTypeRef().GetRef().GetPath(); path != nil {
				appName = paramNameParts[0]
				typeName = path[0]
			} else {
				typeName = paramNameParts[0]
			}
		} else {
			typeName = param.Type.GetTypeRef().GetRef().GetPath()[0]
		}
		parsedType = NewTypeRef(appName, typeName)
		relatedTypes = catalogdiagrams.RecurseivelyGetTypes(appName, map[string]*sysl.Type{typeName: parsedType}, p.Module)
	}
	if _, ok := p.Module.Apps[appName]; ok {
		packageName, _ := GetAppPackageName(p.Module.Apps[appName])
		return p.CreateFileName(catalogdiagrams.GenerateDataModel(appName, relatedTypes), packageName, appName, typeName+"full.svg")
	}
	return ""
}
