package catalog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/anz-bank/sysl/pkg/diagrams"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/cheggaaa/pb/v3"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

// Generator is the contextual object that is used in the markdown generation
type Generator struct {
	FilesToCreate   map[string]string
	Title           string
	Server          bool // Determines whether html refresh header is added
	LiveReload      bool
	Output          string
	PlantumlService string
	OutputFileName  string
	Log             *logrus.Logger
	Fs              afero.Fs
	Module          *sysl.Module
	Format          string
	ProjectTempl    *template.Template // Templ is used to template the Generator struct
	PackageTempl    *template.Template // PackageTempl is passed down to all Packages
}

// NewProject generates a Generator object
func NewProject(
	title, outputDir, plantumlservice, outputType string,
	log *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs) *Generator {

	p := Generator{
		Title:           strings.ReplaceAll(filepath.Base(title), ".sysl", ""),
		Output:          outputDir,
		Log:             log,
		Module:          module,
		PlantumlService: plantumlservice,
		FilesToCreate:   make(map[string]string),
		Fs:              fs,
	}
	var ProjectTemplate, PackageTemplate string
	switch outputType {
	case "html":
		p.OutputFileName = "index.html"
		ProjectTemplate = strings.ReplaceAll(NewProjectTemplate, "README.md", "index.html")
		PackageTemplate = strings.ReplaceAll(NewPackageTemplate, "README.md", "index.html")
		p.Format = "html"
	case "markdown", "md":
		p.OutputFileName = "README.md"
		ProjectTemplate, PackageTemplate = NewProjectTemplate, NewPackageTemplate
		p.Format = "md"
	}
	var err error
	p.ProjectTempl, err = template.New("project").Funcs(p.GetFuncMap()).Parse(ProjectTemplate)
	if err != nil {
		panic(err)
	}
	p.PackageTempl, err = template.New("package").Funcs(p.GetFuncMap()).Parse(PackageTemplate)
	if err != nil {
		panic(err)
	}
	return &p
}

// Run Executes a project and generates markdown and diagrams to a given filesystem
func (p *Generator) Run() {
	packages := ModuleAsPackages(p.Module)
	p.CreateMarkdown(p.ProjectTempl, path.Join(p.Output, p.OutputFileName), p.Module)
	for _, key := range AlphabeticalModules(packages) {
		pkg := packages[key]
		p.CreateMarkdown(p.PackageTempl, path.Join(p.Output, key, p.OutputFileName), pkg)
	}
	if !p.LiveReload {
		var wg sync.WaitGroup
		fmt.Println("Generating diagrams:")
		progress := pb.StartNew(len(p.FilesToCreate))
		//numRoutines := 0
		for key, val := range p.FilesToCreate {
			//go func(key, val string) {
			//wg.Add(1)
			//numRoutines++
			p.Fs.MkdirAll(path.Join(p.Output, path.Dir(key)), os.ModePerm)
			if err := diagrams.OutputPlantuml(path.Join(p.Output, key), p.PlantumlService, val, p.Fs); err != nil {
				panic(err)
			}
			progress.Increment()
			//wg.Done()
			//
			//}(key, val)
			//if numRoutines >= 1 {
			//	wg.Wait()
			//	numRoutines = 0
			//}
		}
		wg.Wait()
		progress.Finish()
	}
}

// GetFuncMap returns the funcs that are used in diagram generation
func (p *Generator) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateIntegrationDiagram":  p.CreateIntegrationDiagram,
		"CreateSequenceDiagram":     p.CreateSequenceDiagram,
		"CreateParamDataModel":      p.CreateParamDataModel,
		"CreateReturnDataModel":     p.CreateReturnDataModel,
		"CreateTypeDiagram":         p.CreateTypeDiagram,
		"GenerateDataModel":         p.GenerateDataModel,
		"CreateQueryParamDataModel": p.CreateQueryParamDataModel,
		"CreatePathParamDataModel":  p.CreatePathParamDataModel,
		"hasPattern":                syslutil.HasPattern,
		"ModuleAsPackages":          ModuleAsPackages,
		"ModulePackageName":         ModulePackageName,
		"AlphabeticalModules":       AlphabeticalModules,
		"AlphabeticalParams":        AlphabeticalParams,
		"AlphabeticalApps":          AlphabeticalApps,
		"AlphabeticalTypes":         AlphabeticalTypes,
		"AlphabeticalEndpoints":     AlphabeticalEndpoints,
		"AppComment":                AppComment,
		"TypeComment":               TypeComment,
		"Attribute":                 Attribute,
		"SanitiseOutputName":        SanitiseOutputName,
	}
}
