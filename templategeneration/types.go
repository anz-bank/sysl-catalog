package templategeneration

import (
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

const (
	diagram_integration = "integration"
	diagram_sequence    = "sequence"
	diagram_datamodel   = "data"
)

// Project is the top level in the hierarchy of markdown generation
type Project struct {
	Title                       string
	Output                      string
	PlantumlService             string
	OutputFileName              string
	RootLevelIntegrationDiagram *Diagram
	Log                         *logrus.Logger
	Packages                    map[string]*Package //Packages are the rows of the top level markdown
	Fs                          afero.Fs
	Module                      *sysl.Module
	PackageModules              map[string]*sysl.Module // PackageModules maps @package attr to all those applications
	ProjectTempl                *template.Template      // Templ is used to template the Project struct
	PackageTempl                *template.Template      // PackageTempl is passed down to all Packages
	EmbededTempl                *template.Template      // This is passed down to all Diagrams
}

// NewProject generates a Project Markdwon object for all a sysl module
func NewProject(inputSyslFileName, output, plantumlservice string, log *logrus.Logger, fs afero.Fs, module *sysl.Module) *Project {
	p := Project{
		Title:           strings.ReplaceAll(filepath.Base(inputSyslFileName), ".sysl", ""),
		Output:          output,
		Fs:              fs,
		Log:             log,
		Module:          module,
		Packages:        map[string]*Package{},
		PackageModules:  map[string]*sysl.Module{},
		PlantumlService: plantumlservice,
		OutputFileName:  pageFilename,
	}
	p.initProject()
	if err := p.RegisterSequenceDiagrams(); err != nil {
		p.Log.Errorf("Error creating parsing sequence diagrams: %v", err)
	}
	return &p
}

// RegisterTemplates registers templates for the project to use
func (p *Project) RegisterTemplates(projectTemplateString, packageTemplateString, embededTemplateString string) error {
	templates, err := LoadMarkdownTemplates(projectTemplateString, packageTemplateString, embededTemplateString)
	if err != nil {
		return err
	}
	p.ProjectTempl, p.PackageTempl, p.EmbededTempl = templates[0], templates[1], templates[2]
	return nil
}

func (p *Project) initProject() {
	for _, key := range alphabeticalApps(p.Module.Apps) {
		app := p.Module.Apps[key]
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
				OutputFile:  pageFilename,
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
func (p *Project) ExecuteTemplateAndDiagrams() {
	if p.EmbededTempl == nil || p.PackageTempl == nil || p.ProjectTempl == nil {
		if err := p.RegisterTemplates(ProjectMarkdownTemplate, PackageMarkdownTemplate, EmbededSvgTemplate); err != nil {
			p.Log.Errorf("Error registering default templates:\n %v", err)
		}
	}
	if err := p.CreateIntegrationDiagrams(); err != nil {
		p.Log.Errorf("Error generating integration diagrams:\n %v", err)
		return
	}

	if err := GenerateMarkdown(p.Output, p.OutputFileName, p, p.ProjectTempl, p.Fs); err != nil {
		p.Log.Errorf("Error generating root markdown:\n %v", err)
		return
	}
	for _, key := range alphabeticalPackage(p.Packages) {
		pkg := p.Packages[key]
		if err := GenerateMarkdown(pkg.OutputDir, pkg.OutputFile, pkg, pkg.Parent.PackageTempl, p.Fs); err != nil {
			p.Log.Errorf("Error generating package markdown:\n %v", err)
			return
		}
		for _, sd := range pkg.SequenceDiagrams {
			if err := GenerateDiagramAndMarkdown(sd); err != nil {
				p.Log.Errorf("Error generating Sequence diagram template and diagrams:\n %v", err)
				return
			}
		}
		for _, intDiagrams := range pkg.IntegrationDiagrams {
			if err := GenerateDiagramAndMarkdown(intDiagrams); err != nil {
				p.Log.Errorf("Error generating Integration diagram template and diagrams:\n %v", err)
				return
			}
		}
		// if _, err := p.CreateDataModelDiagram(); err != nil {
		// 	p.Log.Errorf("Error generating DataModel diagrams:\n %v", err)
		// 	return
		// }
	}
}
