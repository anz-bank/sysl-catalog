// util.go: misc functions to convert/send http requests/sort maps
package catalog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/anz-bank/protoc-gen-sysl/syslpopulate"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
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
				Ref: &sysl.Scope{Appname: &sysl.AppName{
					Part: []string{appName},
				},
					Path: []string{appName, typeName},
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
	app := syslpopulate.NewApplication("")
	app.Endpoints = make(map[string]*sysl.Endpoint)
	app.Endpoints["_"] = syslpopulate.NewEndpoint("_")
	app.Endpoints["_"].Stmt = []*sysl.Statement{}
	for key, _ := range Apps {
		app.Endpoints["_"].Stmt = append(app.Endpoints["_"].Stmt, syslpopulate.NewStringStatement(key))
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

type Attr interface {
	GetAttrs() map[string]*sysl.Attribute
}

func Attribute(a Attr, query string) string {
	if description := a.GetAttrs()[query]; description != nil {
		return description.GetS()
	}
	return ""
}

type Namer interface {
	Attr
	GetName() *sysl.AppName
}

// GetAppPackageName returns the package and app name of any sysl application
func GetAppPackageName(a Namer) (string, string) {
	appName := strings.Join(a.GetName().GetPart(), "")
	packageName := appName
	if attr := a.GetAttrs()["package"]; attr != nil {
		packageName = attr.GetS()
	}
	return packageName, appName
}

func ModuleAsPackages(m *sysl.Module) map[string]*sysl.Module {
	packages := make(map[string]*sysl.Module)
	for _, key := range SortedKeys(m.GetApps()) {
		app := m.GetApps()[key]
		packageName, _ := GetAppPackageName(app)
		if syslutil.HasPattern(app.GetAttrs(), "ignore") {
			continue
		}
		if _, ok := packages[packageName]; !ok {
			packages[packageName] = &sysl.Module{Apps: map[string]*sysl.Application{}}
		}
		packages[packageName].GetApps()[strings.Join(app.GetName().GetPart(), "")] = app
	}
	return packages
}

func ModulePackageName(m *sysl.Module) string {
	for _, key := range SortedKeys(m.GetApps()) {
		app := m.Apps[key]
		packageName, _ := GetAppPackageName(app)
		return packageName
	}
	return ""
}

// Map applies a function to every element in a string slice
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// RetryHTTPRequest retries the given request
func RetryHTTPRequest(url string) ([]byte, error) {
	client := retryablehttp.NewClient()
	client.Logger = nil
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// PlantUMLURL returns a PlantUML url
func PlantUMLURL(plantumlService, contents string) (string, error) {
	encoded, err := diagrams.DeflateAndEncode([]byte(contents))
	if err != nil {
		return "", err
	}
	return fmt.Sprint(plantumlService, "/", "svg", "/", encoded), nil
}

func HttpToFile(url, fileName string, fs afero.Fs) error {
	if err := fs.MkdirAll(path.Dir(fileName), os.ModePerm); err != nil {
		return err
	}
	out, err := RetryHTTPRequest(url)
	if err != nil {
		return err
	}
	if err := afero.WriteFile(fs, fileName, append(out, byte('\n')), os.ModePerm); err != nil {
		return err
	}
	return nil
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

type Typer interface {
	GetType() *sysl.Type
}

// GetAppTypeName returns the appName and typeName of a param
func GetAppTypeName(param Typer) (string, string) {
	var appName, typeName string
	ref := param.GetType().GetTypeRef().GetRef()
	appNameParts := ref.GetAppname().GetPart()
	if len(appNameParts) > 0 {
		typeNameParts := ref.GetPath()
		if typeNameParts != nil {
			appName = appNameParts[0]
			typeName = typeNameParts[0]
		} else {
			typeName = appNameParts[0]
		}
	} else {
		typeName = ref.GetPath()[0]
	}
	return appName, typeName
}
