// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/iancoleman/strcase"

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
	ProjectTitle         string
	RootModule           *sysl.Module
	LiveReload           bool // Add live reload javascript to html
	ImageTags            bool // embedded plantuml img tags, or generated svgs
	DisableCss           bool // used for rendering raw markdown
	DisableImages        bool // used for omitting image creation
	Mermaid              bool
	Format               string // "html" or "markdown" or "" if custom
	CurrentDir           string
	TempDir              string
	Ext                  string
	OutputFileName       string
	PlantumlService      string
	Log                  *logrus.Logger
	Fs                   afero.Fs
	errs                 []error // Any errors that stop from rendering will be output to the browser
	Templates            []*template.Template
	StartTemplateIndex   int
	// All of these are used in markdown generation
	Module    *sysl.Module
	Title     string
	OutputDir string
	Links     map[string]string
}

type SourceCoder interface {
	Attr
	GetSourceContext() *sysl.SourceContext
}

// RootPath appends CurrentDir to output
func (p *Generator) SourcePath(a SourceCoder) string {
	rootDir := rootDirectory(path.Join(p.OutputDir, p.CurrentDir))
	if source_path := Attribute(a, "source_path"); source_path != "" {
		return rootDir + Attribute(a, "source_path")
	}
	return path.Join(rootDir, a.GetSourceContext().File)
}

// NewProject generates a Generator object, fs and outputDir are optional if being used for a web server.
func NewProject(
	titleAndFileName, plantumlService, outputType string,
	log *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs, outputDir string, mermaidEnabled bool) *Generator {
	p := Generator{
		ProjectTitle:    titleAndFileName,
		SourceFileName:  titleAndFileName,
		OutputDir:       outputDir,
		OutputFileName:  typeMaps[strings.ToLower(outputType)],
		Format:          strings.ToLower(outputType),
		Log:             log,
		RootModule:      module,
		PlantumlService: plantumlService,
		FilesToCreate:   make(map[string]string),
		Fs:              fs,
		Ext:             ".svg",
		Mermaid:         mermaidEnabled,
	}
	if module != nil && len(p.ModuleAsMacroPackage(module)) <= 1 {
		p.StartTemplateIndex = 1 // skip the MacroPackageProject
	}
	return p.WithTemplateString(MacroPackageProject, ProjectTemplate, NewPackageTemplate)
}

// WithTemplateFileNames loads template strings into project and package of p respectively
func (p *Generator) WithTemplateString(tmpls ...string) *Generator {
	for i, e := range tmpls {
		tmpl, err := template.New(strconv.Itoa(i)).Funcs(p.GetFuncMap()).Parse(e)
		if err != nil {
			p.Log.Error("Error registering template:", err)
			return nil
		}
		p.Templates = append(p.Templates, tmpl)
	}
	return p
}

func (p *Generator) WithTemplateFs(fs afero.Fs, fileNames ...string) *Generator {
	var tmpls []string
	if len(fileNames) == 0 || fileNames[0] == "" {
		return p
	}
	for _, e := range fileNames {
		bytes, err := afero.ReadFile(fs, e)
		if err != nil {
			p.Log.Error("Error opening template file:", p)
		}
		tmpls = append(tmpls, string(bytes))
	}
	p.Templates = make([]*template.Template, 0, 2)
	p.StartTemplateIndex = 0
	return p.WithTemplateString(tmpls...)
}

func (p *Generator) SetOptions(disableCss, disableImages, imageTags bool, readmeName string) *Generator {
	p.DisableCss = disableCss
	p.DisableImages = disableImages || imageTags
	p.ImageTags = imageTags
	if readmeName != "" {
		p.OutputFileName = readmeName
	}
	return p

}

// Run Executes a project and generates markdown and diagrams to a given filesystem.
func (p *Generator) Run() {
	p.Title = p.ProjectTitle
	fileName := markdownName(p.OutputFileName, path.Base(p.ProjectTitle))
	p.Module = p.RootModule
	if err := p.CreateMarkdown(p.Templates[p.StartTemplateIndex], path.Join(p.OutputDir, fileName), p); err != nil {
		p.Log.Error("Error creating project markdown:", err)
	}
	var wg sync.WaitGroup
	var progress *pb.ProgressBar
	var completedDiagrams int64
	var diagramCreator = func(inMap map[string]string, f func(fs afero.Fs, filename string, data string) error) {
		for fileName, contents := range inMap {
			wg.Add(1)
			go func(fileName, contents string) {
				var err = f(p.Fs, path.Join(p.OutputDir, fileName), contents)
				if err != nil {
					p.Log.Error("Error generating file:", err)
				}
				progress.Increment()
				wg.Done()
				completedDiagrams++
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
	progress.SetCurrent(completedDiagrams)
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

func Last(i interface{}, ind int) bool {
	return ind == len(SortedKeys(i))-1
}

func Remove(s string, old ...string) string {
	for _, e := range old {
		re := regexp.MustCompile(e)
		s = re.ReplaceAllString(s, "")
	}
	return s
}

// GetFuncMap returns the funcs that are used in diagram generation.
func (p *Generator) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateIntegrationDiagram": p.CreateIntegrationDiagram,
		"CreateSequenceDiagram":    p.CreateSequenceDiagram,
		"CreateParamDataModel":     p.CreateParamDataModel,
		"CreateReturnDataModel":    p.CreateReturnDataModel,
		"CreateTypeDiagram":        p.CreateTypeDiagram,
		"GenerateDataModel":        p.GenerateDataModel,
		"GetParamType":             p.GetParamType,
		"GetReturnType":            p.GetReturnType,
		"SourcePath":               p.SourcePath,
		"Packages":                 p.Packages,
		"MacroPackages":            p.MacroPackages,
		"hasPattern":               syslutil.HasPattern,
		"ModuleAsPackages":         p.ModuleAsPackages,
		"ModulePackageName":        ModulePackageName,
		"SortedKeys":               SortedKeys,
		"Attribute":                Attribute,
		"Fields":                   Fields,
		"FieldType":                FieldType,
		"SanitiseOutputName":       SanitiseOutputName,
		"ToLower":                  strings.ToLower,
		"ToCamel":                  strcase.ToCamel,
		"Remove":                   Remove,
		"ToTitle":                  strings.ToTitle,
		"Base":                     filepath.Base,
		"Last":                     Last,
	}
}
