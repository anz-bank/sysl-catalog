// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/cheggaaa/pb/v3"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

var typeMaps = map[string]string{"md": "README.md", "markdown": "README.md", "html": "index.html"}

// Generator is the contextual object that is used in the markdown generation
type Generator struct {
	FilesToCreate        map[string]string
	MermaidFilesToCreate map[string]string
	Title                string
	LiveReload           bool   // Add live reload javascript to html
	ImageTags            bool   // embedded plantuml img tags, or generated svgs
	DisableCss           bool   // used for rendering raw markdown
	DisableImages        bool   // used for omitting image creation
	Format               string // "html" or "markdown" or "" if custom
	OutputDir            string
	OutputFileName       string
	PlantumlService      string
	Log                  *logrus.Logger
	Fs                   afero.Fs
	Module               *sysl.Module
	ProjectTempl         *template.Template // ProjectTempl is used to template the Generator struct
	PackageTempl         *template.Template // PackageTempl is passed down to all Packages
}

// NewProject generates a Generator object, fs and outputDir are optional if being used for a web server.
func NewProject(
	title, plantumlService, outputType string,
	log *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs, outputDir string) *Generator {
	p := Generator{
		Title:           title,
		OutputDir:       outputDir,
		OutputFileName:  typeMaps[strings.ToLower(outputType)],
		Format:          strings.ToLower(outputType),
		Log:             log,
		Module:          module,
		PlantumlService: plantumlService,
		FilesToCreate:   make(map[string]string),
		Fs:              fs,
	}
	return p.WithTemplateString(NewProjectTemplate, NewPackageTemplate)
}

// WithTemplateFileNames loads template strings into project and package of p respectively
func (p *Generator) WithTemplateString(p1, p2 string) *Generator {
	var err error
	if p1 != "" {
		p.ProjectTempl, err = template.New("project").Funcs(p.GetFuncMap()).Parse(p1)
		if err != nil {
			p.Log.Error(err)
			return nil
		}
	}
	if p2 != "" {
		p.PackageTempl, err = template.New("package").Funcs(p.GetFuncMap()).Parse(p2)
		if err != nil {
			p.Log.Error(err)
			return nil
		}
	}
	return p
}

// WithTemplateFileNames loads templates from fs registered in p
func (p *Generator) WithTemplateFiles(p1, p2 afero.File) *Generator {
	var file1, file2 []byte
	var err error
	if p1 != nil {
		file1, err = afero.ReadAll(p1)
		if err != nil {
			p.Log.Error(err)
			return nil
		}
		p.DisableCss = true
		p.Format = ""
	}
	if p2 != nil {
		file2, err = afero.ReadAll(p2)
		if err != nil {
			p.Log.Error(err)
			return nil
		}
		p.DisableCss = true
		p.Format = ""
	}
	return p.WithTemplateString(string(file1), string(file2))
}

func (p *Generator) SetOptions(disableCss, disableImages bool) *Generator {
	p.DisableCss = disableCss
	p.DisableImages = disableImages
	return p

}

// Run Executes a project and generates markdown and diagrams to a given filesystem.
func (p *Generator) Run() {
	m := struct {
		*sysl.Module
		Title string
	}{p.Module, p.Title}
	p.CreateMarkdown(p.ProjectTempl, path.Join(p.OutputDir, p.OutputFileName), m)
	packages := ModuleAsPackages(p.Module)
	for _, key := range AlphabeticalModules(packages) {
		p.CreateMarkdown(p.PackageTempl, path.Join(p.OutputDir, key, p.OutputFileName), packages[key])
	}
	if p.ImageTags || p.DisableImages {
		fmt.Println("Skipping Image creation")
		return
	}
	var wg sync.WaitGroup
	fmt.Println("Generating diagrams:")
	progress := pb.StartNew(len(p.FilesToCreate))
	for key, url := range p.FilesToCreate {
		wg.Add(1)
		go func(key, url string) {
			if err := HttpToFile(url, path.Join(p.OutputDir, key), p.Fs); err != nil {
				p.Log.Error(err)
			}
			progress.Increment()
			wg.Done()
		}(key, url)
	}
	wg.Wait()
	progress.Finish()
}

// GetFuncMap returns the funcs that are used in diagram generation.
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
		"GetParamType":              p.GetParamType,
		"GetReturnType":             p.GetReturnType,
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
		"ToLower":                   strings.ToLower,
		"Base":                      filepath.Base,
	}
}
