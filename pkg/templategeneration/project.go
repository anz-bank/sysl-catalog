package templategeneration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/anz-bank/sysl-catalog/pkg/catalogdiagrams"
	"github.com/anz-bank/sysl/pkg/datamodeldiagram"
	"github.com/anz-bank/sysl/pkg/sequencediagram"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

const (
	diagram_sequence = "sequence"
	refreshHeader    = ``
)

// Project is the top level in the hierarchy of markdown generation
type Project struct {
	Title                          string
	PageFileName                   string // Is README.md for markdown and index.html for html
	Server                         bool   // Determines wether html refresh header is added
	OutputType                     string
	Output                         string
	PlantumlService                string
	OutputFileName                 string
	RootLevelIntegrationDiagramEPA *Diagram
	RootLevelIntegrationDiagram    *Diagram
	Log                            *logrus.Logger
	Packages                       map[string]*Package // Packages are the rows of the top level markdown
	Fs                             afero.Fs
	Module                         *sysl.Module
	PackageModules                 map[string]*sysl.Module // PackageModules maps @package attr to all those applications
	ProjectTempl                   *template.Template      // Templ is used to template the Project struct
	PackageTempl                   *template.Template      // PackageTempl is passed down to all Packages
}

// NewProject generates a Project Markdwon object for all a sysl module
func NewProject(inputSyslFileName, output, plantumlservice string, outputType string, log *logrus.Logger, fs afero.Fs, module *sysl.Module) *Project {
	fileName := "README.md"
	if outputType == "html" {
		fileName = "index.html"
	}
	p := Project{
		Title:           strings.ReplaceAll(filepath.Base(inputSyslFileName), ".sysl", ""),
		Output:          output,
		Fs:              fs,
		Log:             log,
		Module:          module,
		Packages:        map[string]*Package{},
		PackageModules:  map[string]*sysl.Module{},
		PlantumlService: plantumlservice,
		OutputFileName:  fileName,
		OutputType:      outputType,
	}
	p.initProject()
	if err := p.RegisterDiagrams(); err != nil {
		p.Log.Errorf("Error creating parsing sequence diagrams: %v", err)
	}
	return &p
}

// RegisterTemplates registers templates for the project to use
func (p *Project) RegisterTemplates(projectTemplateString, packageTemplateString string) error {
	templates, err := LoadMarkdownTemplates(projectTemplateString, packageTemplateString)
	if err != nil {
		return err
	}
	p.ProjectTempl, p.PackageTempl = templates[0], templates[1]
	return nil
}

// initProject reshuffles apps into "packages"; sort of like "sub modules"
func (p *Project) initProject() {
	for _, key := range AlphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		fmt.Println(p.Title, key)
		if syslutil.HasPattern(app.Attrs, "ignore") || key == p.Title {
			continue
		}
		packageName, _ := GetAppPackageName(app)
		newPackage, ok := p.Packages[packageName]
		if !ok {
			newPackage = &Package{
				Parent:      p,
				PackageName: packageName,
				OutputDir:   path.Join(p.Output, packageName),
				OutputFile:  p.OutputFileName,
			}
		}
		p.Packages[packageName] = newPackage
		if _, ok := p.PackageModules[packageName]; !ok {
			p.PackageModules[packageName] = &sysl.Module{}
			p.PackageModules[packageName].Apps = map[string]*sysl.Application{}
		}
		p.PackageModules[packageName].Apps[strings.Join(app.Name.Part, "")] = app

	}
}

// SetServerMode can be used to add html headers to the html templates for server mode
func (p *Project) SetServerMode() *Project {
	//Add the refresh header is server mode has been enabled
	if err := p.RegisterTemplates(refreshHeader+ProjectHTMLTemplate, refreshHeader+PackageHTMLTemplate); err != nil {
		p.Log.Errorf("Error registering default templates:\n %v", err)
	}
	return p
}

// ExecuteTemplateAndDiagrams generates all documentation of Project with the registered Markdown
func (p *Project) ExecuteTemplateAndDiagrams() {
	var wg sync.WaitGroup
	if p.PackageTempl == nil || p.ProjectTempl == nil {
		if p.OutputType == "html" {
			if err := p.RegisterTemplates(ProjectHTMLTemplate, PackageHTMLTemplate); err != nil {
				p.Log.Errorf("Error registering default templates:\n %v", err)
			}
		} else {
			if err := p.RegisterTemplates(ProjectMarkdownTemplate, PackageMarkdownTemplate); err != nil {
				p.Log.Errorf("Error registering default templates:\n %v", err)
			}
		}
	}
	projectApp, ok := p.Module.Apps[p.Title]
	if !ok {
		projectApp = createProjectApp(p.Module.Apps)
	}
	var err error
	p.RootLevelIntegrationDiagram, err = p.CreateIntegrationDiagrams(p.Title, p.Output, projectApp, false)
	if err != nil {
		p.Log.Errorf("Error generating integration diagrams:\n %v", err)
		return
	}
	p.RootLevelIntegrationDiagramEPA, err = p.CreateIntegrationDiagrams(p.Title, p.Output, projectApp, true)
	if err != nil {
		p.Log.Errorf("Error generating integration diagrams:\n %v", err)
		return
	}
	if err := GenerateMarkdown(p.Output, p.OutputFileName, p, p.ProjectTempl, p.Fs); err != nil {
		p.Log.Errorf("Error generating root markdown:\n %v", err)
		return
	}
	for _, key := range AlphabeticalPackage(p.Packages) {
		pkg := p.Packages[key]
		pkg.Integration, err = p.CreateIntegrationDiagrams(pkg.PackageName, pkg.OutputDir, createProjectApp(p.PackageModules[pkg.PackageName].Apps), false)
		if err != nil {
			p.Log.Errorf("Error generating package int diagram")
		}
		wg.Add(1)
		go func(pk *Package) {
			if err := GenerateMarkdown(pk.OutputDir, pk.OutputFile, pk, pk.Parent.PackageTempl, p.Fs); err != nil {
				p.Log.Errorf("Error generating package markdown:\n %v", err)
				return
			}
			wg.Done()
		}(pkg)
		for _, apps := range pkg.SequenceDiagrams {
			for _, sd := range apps {
				wg.Add(1)
				go func(s *Diagram) {
					if err := s.GenerateDiagramAndMarkdown(); err != nil {
						p.Log.Errorf("Error generating Sequence diagram template and diagrams:\n %v", err)
						return
					}
					wg.Done()
				}(sd)
			}
		}
		for _, data := range pkg.DatabaseModel {
			wg.Add(1)
			go func(s *Diagram) {
				if err := s.GenerateDiagramAndMarkdown(); err != nil {
					p.Log.Errorf("Error generating Sequence diagram template and diagrams:\n %v", err)
					return
				}
				wg.Done()
			}(data)
		}
	}
	wg.Wait()
}

// AlphabeticalRows returns an alphabetically sorted list of packages of any project.
func (p Project) AlphabeticalRows() []*Package {
	packages := make([]*Package, 0, len(p.Packages))
	for _, key := range AlphabeticalPackage(p.Packages) {
		packages = append(packages, p.Packages[key])
	}
	return packages
}

// RegisterDiagrams creates sequence Diagrams from the sysl Module in Project.
func (p Project) RegisterDiagrams() error {
	for _, key := range AlphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]

		packageName, appName := GetAppPackageName(app)
		if syslutil.HasPattern(app.Attrs, "ignore") || key == p.Title {
			p.Log.Infof("Skipping application %s", app.Name)
			continue
		}
		if _, ok := p.Packages[packageName]; !ok {
			p.Packages[packageName] = &Package{Parent: &p}
		}
		if p.Packages[packageName].SequenceDiagrams == nil {
			p.Packages[packageName].SequenceDiagrams = make(map[string][]*Diagram)
			p.Packages[packageName].SequenceDiagrams[appName] = make([]*Diagram, 0, 0)
		}
		if syslutil.HasPattern(app.Attrs, "db") {
			if p.Packages[packageName].DatabaseModel == nil {
				p.Packages[packageName].DatabaseModel = make(map[string]*Diagram)
			}
			p.Packages[packageName].DatabaseModel[appName] = &Diagram{
				Parent:           p.Packages[packageName],
				App:              app,
				DiagramString:    p.GenerateDBDataModel(appName),
				OutputDir:        path.Join(p.Output, packageName),
				OutputFileName__: sanitiseOutputName(appName + "db"),
			}
		}
		if len(app.Endpoints) == 0 {
			continue
		}
		for _, key2 := range AlphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			if syslutil.HasPattern(endpoint.Attrs, "ignore") || key == p.Title {
				p.Log.Infof("Skipping application %s", app.Name)
				continue
			}
			packageD := p.Packages[packageName]
			diagram, err := packageD.SequenceDiagramFromEndpoint(appName, endpoint)
			if err != nil {
				return err
			}

			p.Packages[packageName].SequenceDiagrams[appName] = append(packageD.SequenceDiagrams[appName], diagram)
		}

	}
	return nil
}

// GenerateDBDataModel takes all the types in parentAppName and generates data model diagrams for it
func (p Project) GenerateDBDataModel(parentAppName string) string {
	pl := &datamodelCmd{}
	pl.Project = ""
	p.Fs.MkdirAll(pl.Output, os.ModePerm)
	pl.Direct = true
	pl.ClassFormat = "%(classname)"
	spclass := sequencediagram.ConstructFormatParser("", pl.ClassFormat)
	var stringBuilder strings.Builder
	dataParam := &catalogdiagrams.DataModelParam{}
	dataParam.Mod = p.Module

	v := datamodeldiagram.MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
	vNew := &catalogdiagrams.DataModelView{
		DataModelView: *v,
		TypeMap:       p.Module.Apps[parentAppName].Types,
	}
	return vNew.GenerateDataView(dataParam, parentAppName, nil, p.Module)
}

// GenerateEndpointDataModel generates data model diagrams for a specific type
func (p Project) GenerateEndpointDataModel(parentAppName string, t *sysl.Type) string {
	pl := &datamodelCmd{}
	pl.Project = ""
	p.Fs.MkdirAll(pl.Output, os.ModePerm)
	pl.Direct = true
	pl.ClassFormat = "%(classname)"
	spclass := sequencediagram.ConstructFormatParser("", pl.ClassFormat)
	var stringBuilder strings.Builder
	dataParam := &catalogdiagrams.DataModelParam{}
	dataParam.Mod = p.Module

	v := datamodeldiagram.MakeDataModelView(spclass, dataParam.Mod, &stringBuilder, dataParam.Title, "")
	vNew := &catalogdiagrams.DataModelView{
		DataModelView: *v,
		TypeMap:       map[string]*sysl.Type{},
	}
	return vNew.GenerateDataView(dataParam, parentAppName, t, p.Module)
}
