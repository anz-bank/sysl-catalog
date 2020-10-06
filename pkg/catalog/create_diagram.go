// diagram-creation.go: all the methods attached to the generator object to be used in templating
package catalog

import (
	"os"
	"path"
	"strings"

	"github.com/anz-bank/protoc-gen-sysl/newsysl"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/syslwrapper"
)

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
		appName = JoinAppNameString(endpoint.GetSource())
	}
	return p.RootModule.GetApps()[appName].GetTypes()[typeName]
}

func (p *Generator) ExtractTypeInfo(app *sysl.Application, param Param) (appName, typeName, aliasTypeName string, getRecursive bool) {
	defer func() {
		if err := recover(); err != nil {
			p.Log.Errorf("error creating param data model: %s", err)
		}
	}()
	aliasTypeName = param.GetName()
	appName, typeName = GetAppTypeName(param)
	if aliasTypeName == "" {
		aliasTypeName = typeName
	}
	if appName == "" {
		appName = path.Join(app.GetName().GetPart()...)
	}
	if appName == "primitive" {
		getRecursive = false
	} else {
		getRecursive = true
	}
	return
}

func (p *Generator) ExtractReturnInfo(appname string, stmnt *sysl.Statement, endpoint *sysl.Endpoint) (appName string, typeName string, typeref *sysl.Type, getRecursive bool) {
	var sequence bool
	appName = appname
	ret := stmnt.GetRet()
	if ret == nil {
		return
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
		return
	}
	if syslwrapper.IsPrimitive(typeName) {
		getRecursive = false
		typeref = syslwrapper.MakePrimitive(typeName)
	} else {
		getRecursive = true
	}
	return
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

func (p *Generator) getProjectApp(m *sysl.Module) (*sysl.Application, map[string]string) {
	includedProjects := Filter(
		SortedKeys(m.Apps),
		func(i string) bool {
			return syslutil.HasPattern(m.GetApps()[i].GetAttrs(), "project") &&
				path.Base(m.GetApps()[i].SourceContext.File) == path.Base(p.SourceFileName)
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
		packageName := GetPackageName(p.RootModule, app)
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
		packages[projectEndpoint].GetApps()[GetAppNameString(app)] = app
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
		p.Links = map[string]string{
			"Back": "../" + p.OutputFileName,
		}
		newGenerator := *p
		newGenerator.Module = macroPackage
		err := newGenerator.CreateMarkdown(newGenerator.Templates[1], macroPackageFileName, newGenerator)
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
		packageName, _ := GetAppPackageName(app)
		if packageName == "" {
			packageName = key
		}
		if _, ok := includedProjects[packageName]; len(includedProjects) > 0 && !ok {
			continue
		}
		packageName = GetPackageName(p.RootModule, app)
		if syslutil.HasPattern(app.GetAttrs(), "ignore") || syslutil.HasPattern(app.GetAttrs(), "project") {
			continue
		}
		if _, ok := packages[packageName]; !ok {
			packages[packageName] = &sysl.Module{Apps: map[string]*sysl.Application{}}
		}
		packages[packageName].GetApps()[GetAppNameString(app)] = app
	}
	return packages
}
