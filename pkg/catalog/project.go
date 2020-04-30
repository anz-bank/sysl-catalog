package catalog

import (
	"net/http"
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

// Project is the top level in the hierarchy of markdown generation
type Project struct {
	Title                          string
	PageFileName                   string // Is README.md for markdown and index.html for html
	Server                         bool   // Determines wether html refresh header is added
	LiveReload                     bool
	Output                         string
	PlantumlService                string
	OutputFileName                 string
	RootLevelIntegrationDiagramEPA *Diagram
	RootLevelIntegrationDiagram    *Diagram
	Log                            *logrus.Logger
	Packages                       map[string]*Package // Packages are the rows of the top level markdown
	Fs                             afero.Fs
	Module                         *sysl.Module
	DiagramExt                     string                  //.svg or .html if we're in server mode (we don't send requests
	PackageModules                 map[string]*sysl.Module // PackageModules maps @package attr to all those applications
	ProjectTempl                   *template.Template      // Templ is used to template the Project struct
	PackageTempl                   *template.Template      // PackageTempl is passed down to all Packages
}

// SetOutputFs sets the output filesystem
func (p *Project) SetOutputFs(fs afero.Fs) *Project {
	p.Fs = fs
	return p
}

func (p *Project) SetServerMode() *Project {
	p.Server = true
	p.DiagramExt = ".html"
	return p
}

func (p *Project) EnableLiveReload() *Project {
	p.LiveReload = true
	return p
}

func (p *Project) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Fs == nil {
		p.SetOutputFs(afero.NewMemMapFs())
		p.ExecuteTemplateAndDiagrams()
	}
	request := r.RequestURI
	switch path.Ext((request)) {
	case "":
		request += "index.html"
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	bytes, _ := afero.ReadFile(p.Fs, request)
	file := string(bytes)
	if p.LiveReload {
		file = strings.ReplaceAll(string(bytes), "<body>", `<body><script src="/livereload.js?port=6900&mindelay=10&v=2" data-no-instant defer></script>`)
	}
	w.Write([]byte(file))
}

// NewProject generates a Project Markdwon object for all a sysl module
func NewProject(title, outputDir, plantumlservice string, outputType string, log *logrus.Logger, module *sysl.Module) *Project {
	var ProjectTemplate, PackageTemplate string
	p := Project{
		Title:           strings.ReplaceAll(filepath.Base(title), ".sysl", ""),
		Output:          outputDir,
		Log:             log,
		Module:          module,
		DiagramExt:      ".svg",
		Packages:        map[string]*Package{},
		PackageModules:  map[string]*sysl.Module{},
		PlantumlService: plantumlservice,
	}

	switch outputType {
	case "html":
		p.OutputFileName = "index.html"
		ProjectTemplate, PackageTemplate = ProjectHTMLTemplate, PackageHTMLTemplate
	case "markdown", "md":
		p.OutputFileName = "README.md"
		ProjectTemplate, PackageTemplate = ProjectMarkdownTemplate, PackageMarkdownTemplate
	}
	if err := p.RegisterTemplates(ProjectTemplate, PackageTemplate); err != nil {
		p.Log.Errorf("Error registering default templates:\n %v", err)
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

func (p *Project) Update(m *sysl.Module) *Project {
	p.Fs = afero.NewMemMapFs()
	p.Packages = make(map[string]*Package)
	p.PackageModules = make(map[string]*sysl.Module)
	p.Module = m
	p.ExecuteTemplateAndDiagrams()
	return p
}

// initProject reshuffles apps into "packages"; sort of like "sub modules"
func (p *Project) initProject() {
	for _, key := range AlphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
		p.Log.Info(p.Title, key)
		if syslutil.HasPattern(app.Attrs, "ignore") {
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

// ExecuteTemplateAndDiagrams generates all documentation of Project with the registered Markdown
func (p *Project) ExecuteTemplateAndDiagrams() *Project {
	var wg sync.WaitGroup // Make diagram generation concurrent
	var err error
	p.initProject()
	if err := p.RegisterDiagrams(); err != nil {
		p.Log.Errorf("Error creating parsing sequence diagrams: %v", err)
	}
	projectApp := createProjectApp(p.Module.Apps)

	p.RootLevelIntegrationDiagram, err = p.CreateIntegrationDiagrams(p.Title, p.Output, projectApp, false)
	if err != nil {
		p.Log.Errorf("Error generating integration diagrams:\n %v", err)
		return p
	}
	p.RootLevelIntegrationDiagramEPA, err = p.CreateIntegrationDiagrams(p.Title, p.Output, projectApp, true)
	if err != nil {
		p.Log.Errorf("Error generating integration diagrams:\n %v", err)
		return p
	}
	if err := GenerateMarkdown(p.Output, p.OutputFileName, p, p.ProjectTempl, p.Fs); err != nil {
		p.Log.Errorf("Error generating root markdown:\n %v", err)
		return p
	}
	for _, key := range AlphabeticalPackage(p.Packages) {
		pkg := p.Packages[key]
		pkg.Integration, err = p.CreateIntegrationDiagrams(pkg.PackageName, pkg.OutputDir, createProjectApp(p.PackageModules[pkg.PackageName].Apps), false)
		if err != nil {
			p.Log.Errorf("Error generating package int diagram")
		}
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
		defer func() {
			if err := GenerateMarkdown(pkg.OutputDir, pkg.OutputFile, pkg, pkg.Parent.PackageTempl, p.Fs); err != nil {
				p.Log.Errorf("Error generating package markdown:\n %v", err)
			}
		}()

	}
	wg.Wait()
	return p
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
		if syslutil.HasPattern(app.Attrs, "ignore") {
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
				Parent:                p.Packages[packageName],
				App:                   app,
				PlantUMLDiagramString: p.GenerateDBDataModel(appName),
				OutputDir:             path.Join(p.Output, packageName),
				OutputFileName__:      sanitiseOutputName(appName+"db") + p.DiagramExt,
			}
		}
		if len(app.Endpoints) == 0 {
			continue
		}
		for _, key2 := range AlphabeticalEndpoints(app.Endpoints) {
			endpoint := app.Endpoints[key2]
			if syslutil.HasPattern(endpoint.Attrs, "ignore") {
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
	//p.Fs.MkdirAll(pl.Output, os.ModePerm)
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
