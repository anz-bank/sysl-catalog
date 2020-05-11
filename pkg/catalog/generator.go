// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/pkg/errors"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/cheggaaa/pb/v3"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/joshcarp/mermaid-go/mermaid"
	"github.com/spf13/afero"
)

var typeMaps = map[string]string{"md": "README.md", "markdown": "README.md", "html": "index.html"}

// Generator is the contextual object that is used in the markdown generation
type Generator struct {
	FilesToCreate        map[string]string
	MermaidFilesToCreate map[string]string
	SourceFileName       string
	Title                string
	LiveReload           bool // Add live reload javascript to html
	ImageTags            bool // embedded plantuml img tags, or generated svgs
	DisableCss           bool // used for rendering raw markdown
	DisableImages        bool // used for omitting image creation
	Mermaid              bool
	Format               string // "html" or "markdown" or "" if custom
	Ext                  string
	OutputDir            string
	OutputFileName       string
	PlantumlService      string
	Log                  *logrus.Logger
	Fs                   afero.Fs
	Module               *sysl.Module
	errs                 []error // Any errors that stop from rendering will be output to the browser
	ProjectTempl         *template.Template
	PackageTempl         *template.Template
}

// NewProject generates a Generator object, fs and outputDir are optional if being used for a web server.
func NewProject(
	titleAndFileName, plantumlService, outputType string,
	log *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs, outputDir string, mermaidEnabled bool) *Generator {
	p := Generator{
		Title:           titleAndFileName,
		SourceFileName:  titleAndFileName,
		OutputDir:       outputDir,
		OutputFileName:  typeMaps[strings.ToLower(outputType)],
		Format:          strings.ToLower(outputType),
		Log:             log,
		Module:          module,
		PlantumlService: plantumlService,
		FilesToCreate:   make(map[string]string),
		Fs:              fs,
		Ext:             ".svg",
		Mermaid:         mermaidEnabled,
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
func (p *Generator) GetRows(module *sysl.Module) []string {
	if packages := p.ModuleAsPackages2(module); len(packages) > 1 {
		return SortedKeys(packages)
	}
	packages := p.ModuleAsPackages(module)
	return SortedKeys(packages)
}

// Run Executes a project and generates markdown and diagrams to a given filesystem.
func (p *Generator) Run() {
	type mWrap struct {
		*sysl.Module
		Title string
		Links map[string]string
	}
	m := mWrap{Module: p.Module, Title: p.Title}
	if err := p.CreateMarkdown(p.ProjectTempl, path.Join(p.OutputDir, p.OutputFileName), m); err != nil {
		p.Log.Error(err)
	}
	packages := p.ModuleAsPackages2(p.Module)
	if len(packages) <= 1 {
		packages := p.ModuleAsPackages(p.Module)
		for _, key2 := range SortedKeys(packages) {
			pkg := packages[key2]
			fullOutputName := path.Join(p.OutputDir, key2, p.OutputFileName)
			if err := p.CreateMarkdown(p.PackageTempl, fullOutputName, pkg); err != nil {
				p.Log.Error(errors.Wrap(err, "error in generating "+fullOutputName))
			}
		}
	} else {
		for _, key := range SortedKeys(packages) {
			moduleMap := packages[key]
			module := createModuleFromSlices(p.Module, SortedKeys(moduleMap))
			subpackages := p.ModuleAsPackages(module)
			if err := p.CreateMarkdown(p.ProjectTempl, path.Join(p.OutputDir, key, p.OutputFileName), mWrap{Module: module, Title: key, Links: map[string]string{"Back": "../" + p.OutputFileName}}); err != nil {
				p.Log.Error(err)
			}
			for _, key2 := range SortedKeys(subpackages) {
				pkg := subpackages[key2]
				fullOutputName := path.Join(p.OutputDir, key, key2, p.OutputFileName)
				if err := p.CreateMarkdown(p.PackageTempl, fullOutputName, pkg); err != nil {
					p.Log.Error(errors.Wrap(err, "error in generating "+fullOutputName))
				}
			}
		}
	}

	if p.ImageTags || p.DisableImages {
		logrus.Info("Skipping Image creation")
		return
	}
	var wg sync.WaitGroup
	fmt.Println("Generating diagrams:")
	progress := pb.StartNew(len(p.FilesToCreate) + len(p.MermaidFilesToCreate))
	for fileName, url := range p.FilesToCreate {
		wg.Add(1)
		go func(fileName, url string) {
			if err := HttpToFile(url, path.Join(p.OutputDir, fileName), p.Fs); err != nil {
				p.Log.Error(err)
			}
			progress.Increment()
			wg.Done()
		}(fileName, url)
	}
	for fileName, contents := range p.MermaidFilesToCreate {
		wg.Add(1)
		go func(fileName, contents string) {
			mermaidSvg := mermaid.Execute(contents)
			var err = afero.WriteFile(p.Fs, fileName, []byte(mermaidSvg+"\n"), os.ModePerm)
			if err != nil {
				p.Log.Error(err)
			}
			progress.Increment()
			wg.Done()
		}(fileName, contents)
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
		"GetRows":                   p.GetRows,
		"GetReturnType":             p.GetReturnType,
		"hasPattern":                syslutil.HasPattern,
		"ModuleAsPackages":          p.ModuleAsPackages,
		"ModulePackageName":         ModulePackageName,
		"SortedKeys":                SortedKeys,
		"Attribute":                 Attribute,
		"SanitiseOutputName":        SanitiseOutputName,
		"ToLower":                   strings.ToLower,
		"Base":                      filepath.Base,
	}
}
