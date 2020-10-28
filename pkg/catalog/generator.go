// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/joshcarp/gop/gop"

	"github.com/Masterminds/sprig"

	"github.com/iancoleman/strcase"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/syslwrapper"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

var outputFileNames = map[string]string{
	"md":       "README.md",
	"markdown": "README.md",
	"html":     "index.html",
}

// Generator is the contextual object that is used in the markdown generation
type Generator struct {
	RootModule           *sysl.Module
	FilesToCreate        map[string]string
	MermaidFilesToCreate map[string]string
	RedocFilesToCreate   map[string]string
	GeneratedFiles       map[string][]byte
	SourceFileName       string
	ProjectTitle         string
	ImageDest            string // Output all images into this folder is set
	Format               string // "html" or "markdown" or "" if custom
	OutputFileName       string
	PlantumlService      string
	Templates            []*template.Template
	Redoc                *template.Template
	StartTemplateIndex   int
	FilterPackage        []string // Filter these regex terms out of packagenames

	Retriever      gop.Retriever
	CustomTemplate bool
	LiveReload     bool // Add live reload javascript to html
	DisableCss     bool // used for rendering raw markdown

	Log  *logrus.Logger
	Fs   afero.Fs
	errs []error // Any errors that stop from rendering will be output to the browser

	// All of these are used in markdown generation
	Module     *sysl.Module
	CurrentDir string
	TempDir    string
	Title      string
	OutputDir  string
	Links      map[string]string
	Server     bool

	Mapper *syslwrapper.AppMapper

	BasePath string // for using on another endpoint that isn't '/'
}

// SourcePath appends CurrentDir to output
func (p *Generator) SourcePath(a SourceCoder) string {
	//FIXME: handle source_path attr
	// if source_path := Attribute(a, "source_path"); source_path != "" {
	// 	return rootDirectory(path.Join(p.OutputDir, p.CurrentDir)) + Attribute(a, "source_path")
	// }
	// sourcePath, err := handleSourceURL(a.GetSourceContext().File)
	str := BuildSpecURL(a.GetSourceContext().GetFile(), a.GetSourceContext().GetVersion())
	return str
}

// NewProjectFromJson generates a generator object with a json byte input (of a sysl module) instead of a sysl module
func NewProjectFromJson(
	titleAndFileName, plantumlService, outputType string,
	logger *logrus.Logger,
	module []byte,
	fs afero.Fs, outputDir string) *Generator {
	m := &sysl.Module{}
	if err := UnmarshallJson(module, m); err != nil {
		logger.Error("Error unmarshalling data")
		return nil
	}
	return NewProject(titleAndFileName, plantumlService, outputType, logger, m, fs, outputDir)

}

// NewProject generates a Generator object, fs and outputDir are optional if being used for a web server.
func NewProject(
	titleAndFileName, plantumlService, outputType string,
	logger *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs, outputDir string) *Generator {
	p := Generator{
		ProjectTitle:         titleAndFileName,
		SourceFileName:       titleAndFileName,
		OutputDir:            outputDir,
		OutputFileName:       outputFileNames[strings.ToLower(outputType)],
		Format:               strings.ToLower(outputType),
		Log:                  logger,
		RootModule:           module,
		PlantumlService:      plantumlService,
		FilesToCreate:        make(map[string]string),
		GeneratedFiles:       make(map[string][]byte),
		RedocFilesToCreate:   make(map[string]string),
		MermaidFilesToCreate: make(map[string]string),
		Fs:                   fs,
		Redoc:                template.Must(template.New("redoc").Parse(RedocPage)),
	}
	if module != nil && len(p.ModuleAsMacroPackage(module)) <= 1 {
		p.StartTemplateIndex = 1 // skip the MacroPackageProject
	}
	return p.WithTemplateString(MacroPackageProjectMermaid, ProjectTemplateMermaid, NewPackageTemplateMermaid)
}

func (p *Generator) WithRetriever(retr gop.Retriever) *Generator {
	p.Retriever = retr
	return p
}

func (p *Generator) WithRemoteTemplateString(remoteResources ...string) *Generator {
	var tmpls []string
	for _, r := range remoteResources {
		c, _, err := p.Retriever.Retrieve(r)
		if err != nil {
			panic(err)
		}
		tmpls = append(tmpls, string(c))
	}
	return p.WithTemplateString(tmpls...)
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
func (p *Generator) AutomaticTemplates(fs afero.Fs, fileNames ...string) *Generator {
	for _, f := range fileNames {
		switch f {
		case "plantuml":
			return p.WithTemplateString(ProjectTemplate, MacroPackageProject, NewPackageTemplate)
		case "mermaid":
			return p.WithTemplateString(ProjectTemplateMermaid, MacroPackageProjectMermaid, NewPackageTemplateMermaid)
		default:
			if strings.Contains(f, "@") {
				return p.WithRemoteTemplateString(fileNames...)
			}
			return p.WithTemplateFs(fs, fileNames...)
		}
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
			os.Exit(1)
		}
		tmpls = append(tmpls, string(bytes))
	}
	p.Templates = make([]*template.Template, 0, 2)
	p.StartTemplateIndex = 0
	p.CustomTemplate = true
	return p.WithTemplateString(tmpls...)
}

func (p *Generator) SetOptions(
	disableCss bool,
	readmeName, ImageDest string) *Generator {
	p.DisableCss = disableCss
	p.ImageDest = ImageDest
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
	if p.Module != nil {
		p.Mapper = syslwrapper.MakeAppMapper(p.Module)
		p.Mapper.IndexTypes()
		p.Mapper.ConvertTypes()
	}
	if err := p.CreateMarkdown(p.Templates[p.StartTemplateIndex], path.Join(p.OutputDir, fileName), p); err != nil {
		p.Log.Error("Error creating project markdown:", err)
	}
}

// GetFuncMap returns the funcs that are used in diagram generation.
func (p *Generator) GetFuncMap() template.FuncMap {
	f := template.FuncMap{

		/* Mermaid Diagram functions */
		"IntegrationMermaid":     p.IntegrationMermaid,
		"SequenceMermaid":        p.SequenceMermaid,
		"DataModelReturnMermaid": p.DataModelReturnMermaid,
		"DataModelMermaid":       p.DataModelMermaid,
		"DataModelAliasMermaid":  p.DataModelAliasMermaid,
		"DataModelAppMermaid":    p.DataModelAppMermaid,

		// Datamodel table functions
		"DataModelReturnTable": p.DataModelReturnTable,
		"DataModelAliasTable":  p.DataModelAliasTable,
		"DataModelTable":       p.DataModelTable,

		/* Plantuml Diagram functions */

		"IntegrationPlantuml":     p.IntegrationPlantuml,
		"SequencePlantuml":        p.SequencePlantuml,
		"DataModelReturnPlantuml": p.DataModelReturnPlantuml,
		"DataModelParamPlantuml":  p.DataModelParamPlantuml,
		"DataModelAppPlantuml":    p.DataModelAppPlantuml,
		"DataModelPlantuml":       p.DataModelPlantuml,
		"DataModelAliasPlantuml":  p.DataModelAliasPlantuml,

		/* Redoc Functions */
		"CreateRedoc": p.CreateRedoc,

		/* Utility functions */
		"GetParamType":       p.GetParamType,
		"GetReturnType":      p.GetReturnType,
		"SourcePath":         p.SourcePath,
		"Packages":           p.Packages,
		"MacroPackages":      p.MacroPackages,
		"hasPattern":         syslutil.HasPattern,
		"ModuleAsPackages":   p.ModuleAsPackages,
		"ModulePackageName":  ModulePackageName,
		"ModuleNamespace":    ModuleNamespace,
		"SortedKeys":         SortedKeys,
		"Attribute":          Attribute,
		"ServiceMetadata":    ServiceMetadata,
		"Fields":             Fields,
		"FieldType":          FieldType,
		"SanitiseOutputName": SanitiseOutputName,
		"SimpleName":         SimpleName,
		"ToLower":            strings.ToLower,
		"ToCamel":            strcase.ToCamel,
		"Remove":             Remove,
		"ToTitle":            strings.ToTitle,
		"Base":               filepath.Base,
		"Last":               Last,
	}
	for name, function := range sprig.FuncMap() {
		f[name] = function
	}
	return f
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
