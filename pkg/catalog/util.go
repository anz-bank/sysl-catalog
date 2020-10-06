// util.go: misc functions to convert/send http requests/sort maps
package catalog

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/anz-bank/protoc-gen-sysl/newsysl"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/sirupsen/logrus"
)

// SanitiseOutputName removes characters so that the string can be used as a hyperlink.
func SanitiseOutputName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "/", "")
}

func SortedKeys(m interface{}) []string {
	keys := reflect.ValueOf(m).MapKeys()
	ret := make([]string, 0, len(keys))
	for _, v := range keys {
		ret = append(ret, v.String())
	}
	sort.Strings(ret)
	return ret
}

// NewTypeRef returns a type reference, needed to correctly generate data model diagrams
func NewTypeRef(appName, typeName string) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Ref: &sysl.Scope{
					Appname: &sysl.AppName{Part: []string{appName}},
					Path:    []string{typeName},
				},
			},
		},
	}
}

// TernaryOperator returns the first element if bool is true and the second element is false
func TernaryOperator(condition bool, i ...interface{}) interface{} {
	if condition {
		return i[0]
	}
	return i[1]
}

// createProjectApp returns a "project" app used to make integration diagrams for any "sub module" apps
func createProjectApp(Apps map[string]*sysl.Application) *sysl.Application {
	app := newsysl.Application("")
	app.Endpoints = make(map[string]*sysl.Endpoint)
	app.Endpoints["_"] = newsysl.Endpoint("_")
	app.Endpoints["_"].Stmt = []*sysl.Statement{}
	for key := range Apps {
		app.Endpoints["_"].Stmt = append(app.Endpoints["_"].Stmt, newsysl.StringStatement(key))
	}
	if app.Attrs == nil {
		app.Attrs = make(map[string]*sysl.Attribute)
	}
	if _, ok := app.Attrs["appfmt"]; !ok {
		app.Attrs["appfmt"] = &sysl.Attribute{
			Attribute: &sysl.Attribute_S{S: "%(appname)"},
		}
	}
	return app
}

func Attribute(a Attr, query string) string {
	if description := a.GetAttrs()[query]; description != nil {
		return description.GetS()
	}
	return ""
}

func ServiceMetadata(a Attr) string {
	queries := []string{
		"Repo.URL",
		"Owner.Email",
		"Owner.Slack",
		"Server.Prod.URL",
		"Server.UAT.URL",
		"Lifecycle",
	}
	queryMap := make(map[string]string)
	for _, q := range queries {
		queryMap[strings.ToLower(q)] = ""
	}
	for attrName := range a.GetAttrs() {
		q := strings.ToLower(attrName)
		if _, exists := queryMap[q]; exists {
			queryMap[q] = Attribute(a, attrName)
		}
	}

	metadata := strings.Builder{}
	for _, q := range queries {
		if val := queryMap[strings.ToLower(q)]; val != "" {
			metadata.WriteString(fmt.Sprintf("%s: %s\n\n", q, val))
		}
	}
	return metadata.String()
}

func Fields(t *sysl.Type) map[string]*sysl.Type {
	if tuple := t.GetTuple(); tuple != nil {
		return tuple.GetAttrDefs()
	}
	return nil
}

func FieldType(t *sysl.Type) string {
	typeName, typeDetail := syslutil.GetTypeDetail(t)
	if typeName == "primitive" {
		return strings.ToLower(typeDetail)
	}
	if typeName == "sequence" {
		return "sequence of " + typeDetail
	}
	if typeName == "type_ref" {
		return strings.Join(t.GetTypeRef().GetRef().GetPath(), ".")
	}
	if typeName != "" {
		return typeName
	}
	return typeDetail
}

// GetAppNameString returns an app's name as a string, with the namespace joined on "::".
func GetAppNameString(a Namer) string {
	return JoinAppNameString(a.GetName())
}

// JoinAppNameString transforms an AppName to a string, with the namespace joined on "::".
func JoinAppNameString(an *sysl.AppName) string {
	return strings.Join(an.GetPart(), " :: ")
}

// GetAppPackageName returns the package and app name of any sysl application
func GetAppPackageName(a Namer) (string, string) {
	appName := GetAppNameString(a)
	packageName := appName
	if attr := a.GetAttrs()["package"]; attr != nil {
		packageName = attr.GetS()
	}
	return packageName, appName
}

func GetPackageName(m *sysl.Module, a Namer) string {
	packageName, _ := GetAppPackageName(a)
	pkg := m.Apps[packageName]
	if attr := pkg.GetAttrs()["package_alias"]; attr != nil {
		return attr.GetS()
	}
	return packageName

}

func ModulePackageName(m *sysl.Module) string {
	for _, key := range SortedKeys(m.GetApps()) {
		app := m.Apps[key]
		return GetPackageName(m, app)
	}
	return ""
}

// Map applies a function to every element in a string slice
func Filter(vs []string, f func(string) bool) []string {
	vsm := make([]string, 0, len(vs))
	for _, v := range vs {
		if f(v) {
			vsm = append(vsm, v)
		}
	}
	return vsm
}

// PlantUMLURL returns a PlantUML url
func PlantUMLURL(plantumlService, contents string) string {
	encoded, _ := diagrams.DeflateAndEncode([]byte(contents))
	return fmt.Sprint(plantumlService, "/", "svg", "/~1", encoded)
}

// CreateSequenceDiagram creates an sequence diagram and returns the sequence diagram string and any errors
func CreateSequenceDiagram(m *sysl.Module, call string) (string, error) {
	l := &cmdutils.Labeler{}
	p := &sequencediagram.SequenceDiagParam{
		AppLabeler:      l,
		EndpointLabeler: l,
		Endpoints:       []string{call},
		Title:           call,
	}
	return sequencediagram.GenerateSequenceDiag(m, p, logrus.New())
}

// GetAppTypeName takes a Sysl Type and returns the appName and typeName of a param
// If the type is a primitive, the appName returned is "primitive"
func GetAppTypeName(param Typer) (appName string, typeName string) {
	ref := param.GetType().GetTypeRef().GetRef()
	appNameParts := ref.GetAppname().GetPart()
	if a, b := syslutil.GetTypeDetail(param.GetType()); a == "primitive" {
		return a, b
	}
	if len(appNameParts) > 0 {
		typeNameParts := ref.GetPath()
		if typeNameParts != nil {
			appName = JoinAppNameString(ref.GetAppname())
			typeName = typeNameParts[0]
		} else {
			typeName = appNameParts[0]
		}
	} else {
		typeName = ref.GetPath()[0]
	}
	return appName, typeName
}

// Unmarshall Json unmarshalls json bytes into a sysl module
func UnmarshallJson(b []byte, m *sysl.Module) error {
	if m == nil {
		return fmt.Errorf("module is nil")
	}
	ma := protojson.UnmarshalOptions{}
	err := ma.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return err
}
