// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"fmt"
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
	CurrentDir           string
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

func (p *Generator) SetOptions(disableCss, disableImages bool, readmeName string) *Generator {
	p.DisableCss = disableCss
	p.DisableImages = disableImages
	if readmeName != "" {
		p.OutputFileName = readmeName
	}
	return p

}

// GetRows returns a slice of rows that should be output on the index pages of the markdown
func (p *Generator) GetRows(module *sysl.Module) []string {
	if packages := p.ModuleAsMacroPackage(module); len(packages) > 1 {
		keys := SortedKeys(packages)
		return keys
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
	fileName := markdownName(p.OutputFileName, strings.ReplaceAll(path.Base(p.Title), ".sysl", ""))
	if err := p.CreateMarkdown(p.ProjectTempl, path.Join(p.OutputDir, fileName), m); err != nil {
		p.Log.Error(err)
	}
	macroPackages := p.ModuleAsMacroPackage(p.Module)
	var packages map[string]*sysl.Module
	var macroPackageName string
	// We either execute this function when we're iterating through the simple packages
	// or if there are "macroPackages" defined on the ~project app (their endpoints)
	packageFunc := func() {
		for _, packageName := range SortedKeys(packages) {
			pkg := packages[packageName]
			p.CurrentDir = path.Join(macroPackageName, packageName)
			fileName := markdownName(p.OutputFileName, packageName)
			fullOutputName := path.Join(p.OutputDir, p.CurrentDir, fileName)
			if err := p.CreateMarkdown(p.PackageTempl, fullOutputName, pkg); err != nil {
				p.Log.Error(errors.Wrap(err, "error in generating "+fullOutputName))
			}
		}
	}
	switch len(macroPackages) {
	case 0, 1:
		packages = p.ModuleAsPackages(p.Module)
		packageFunc()
	default:
		for _, key := range SortedKeys(macroPackages) {
			macroPackageName = key
			moduleMap := macroPackages[macroPackageName]
			module := createModuleFromSlices(p.Module, SortedKeys(moduleMap))
			packages = p.ModuleAsPackages(module)
			fileName := markdownName(p.OutputFileName, macroPackageName)
			macroPackageFileName := path.Join(p.OutputDir, macroPackageName, fileName)
			p.CurrentDir = macroPackageName
			m := mWrap{Module: module, Title: macroPackageName, Links: map[string]string{"Back": "../" + p.OutputFileName}}
			err := p.CreateMarkdown(p.ProjectTempl, macroPackageFileName, m)
			if err != nil {
				p.Log.Error(err)
			}
			packageFunc()
		}
	}
	var wg sync.WaitGroup
	var progress *pb.ProgressBar
	var done int64
	var diagramCreator = func(inMap map[string]string, f func(fs afero.Fs, filename string, data string) error) {
		for fileName, contents := range inMap {
			wg.Add(1)
			go func(fileName, contents string) {
				var err = f(p.Fs, path.Join(p.OutputDir, fileName), contents)
				if err != nil {
					p.Log.Error(err)
				}
				progress.Increment()
				wg.Done()
				done++
			}(fileName, contents)
		}
	}

	if p.Mermaid {
		progress = pb.StartNew(len(p.MermaidFilesToCreate))
		fmt.Println("Generating Mermaid diagrams:")
		diagramCreator(p.MermaidFilesToCreate, GenerateAndWriteMermaidDiagram)
	}
	if p.ImageTags || p.DisableImages {
		logrus.Info("Skipping Image creation")
		return
	}
	progress = pb.StartNew(len(p.FilesToCreate) + len(p.MermaidFilesToCreate))
	progress.SetCurrent(done)
	fmt.Println("Generating diagrams:")
	diagramCreator(p.FilesToCreate, HttpToFile)

	wg.Wait()
	progress.Finish()
}

func markdownName(s, candidate string) string {
	if strings.Contains(s, "{{.Title}}") {
		candidate = SanitiseOutputName(candidate)
		return strings.ReplaceAll(s, "{{.Title}}", candidate)
	}
	return s
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
