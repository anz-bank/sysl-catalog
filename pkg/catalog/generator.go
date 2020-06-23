// generator.go: struct for converting sysl modules to documentation (Generator)
package catalog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/iancoleman/strcase"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
)

var outputFileNames = map[string]string{"md": "README.md", "markdown": "README.md", "html": "index.html"}

// Generator is the contextual object that is used in the markdown generation
type Generator struct {
	RootModule           *sysl.Module
	FilesToCreate        map[string]string
	MermaidFilesToCreate map[string]string
	RedocFilesToCreate   map[string]string
	SourceFileName       string
	ProjectTitle         string
	ImageDest            string // Output all images into this folder is set
	Format               string // "html" or "markdown" or "" if custom
	Ext                  string
	OutputFileName       string
	PlantumlService      string
	Templates            []*template.Template
	StartTemplateIndex   int

	LiveReload    bool // Add live reload javascript to html
	ImageTags     bool // embedded plantuml img tags, or generated svgs
	DisableCss    bool // used for rendering raw markdown
	DisableImages bool // used for omitting image creation
	Mermaid       bool
	Redoc         bool // used for generating redoc for openapi specs

	Log  *logrus.Logger
	Fs   afero.Fs
	errs []error // Any errors that stop from rendering will be output to the browser

	// All of these are used in markdown generation
	CurrentDir string
	TempDir    string
	Module     *sysl.Module
	Title      string
	OutputDir  string
	Links      map[string]string
	Server     bool
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
	logger *logrus.Logger,
	module *sysl.Module,
	fs afero.Fs, outputDir string) *Generator {
	p := Generator{
		ProjectTitle:       titleAndFileName,
		SourceFileName:     titleAndFileName,
		OutputDir:          outputDir,
		OutputFileName:     outputFileNames[strings.ToLower(outputType)],
		Format:             strings.ToLower(outputType),
		Log:                logger,
		RootModule:         module,
		PlantumlService:    plantumlService,
		FilesToCreate:      make(map[string]string),
		RedocFilesToCreate: make(map[string]string),
		Fs:                 fs,
		Ext:                ".svg",
	}
	//if strings.Contains(p.PlantumlService, ".jar") {
	//	_, err := os.Open(p.PlantumlService)
	//	if err != nil {
	//		p.Log.Error("Error adding plantumlenv:", err)
	//	}
	//}
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
			os.Exit(1)
		}
		tmpls = append(tmpls, string(bytes))
	}
	p.Templates = make([]*template.Template, 0, 2)
	p.StartTemplateIndex = 0
	return p.WithTemplateString(tmpls...)
}

func (p *Generator) SetOptions(
	disableCss, disableImages, imageTags, redoc, mermaidEnabled bool,
	readmeName, ImageDest string) *Generator {
	p.Redoc = redoc
	p.DisableCss = disableCss
	p.DisableImages = disableImages || imageTags
	p.Mermaid = mermaidEnabled
	p.ImageTags = imageTags
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
	if err := p.CreateMarkdown(p.Templates[p.StartTemplateIndex], path.Join(p.OutputDir, fileName), p); err != nil {
		p.Log.Error("Error creating project markdown:", err)
	}

	var progress *pb.ProgressBar
	defer func() {
		if progress != nil {
			progress.Finish()
			fmt.Printf("The generated files are output to folder `%s`\n", p.OutputDir)
		}
	}()

	var wg sync.WaitGroup
	var diagramCreator = func(inMap map[string]string, f func(fs afero.Fs, filename string, data string) error, progress *pb.ProgressBar) {
		for fileName, contents := range inMap {
			wg.Add(1)
			go func(fileName, contents string) {
				var err = f(p.Fs, fileName, contents)
				if err != nil {
					p.Log.Error("Error generating file:", err)
					os.Exit(1)
				}
				if progress != nil {
					progress.Increment()
				}
				wg.Done()
			}(fileName, contents)
		}
	}
	if p.Mermaid {
		progress = pb.Full.Start(len(p.MermaidFilesToCreate))
		diagramCreator(p.MermaidFilesToCreate, GenerateAndWriteMermaidDiagram, progress)
	}
	if p.Redoc {
		progress = pb.Full.Start(len(p.RedocFilesToCreate))
		diagramCreator(p.RedocFilesToCreate, GenerateAndWriteRedoc, progress)
	}
	if (p.ImageTags || p.DisableImages) && !p.Redoc {
		logrus.Info("Skipping Image creation")
		return
	}
	if strings.Contains(p.PlantumlService, ".jar") {
		if !p.Server {
			diagramCreator(p.FilesToCreate, p.PUMLFile, progress)
			start := time.Now()
			if err := PlantUMLJava(p.PlantumlService, p.OutputDir); err != nil {
				p.Log.Error(err)
			}
			elapsed := time.Since(start)
			fmt.Println("Generating took ", elapsed)
		}
	} else {
		progress = pb.Full.Start(len(p.FilesToCreate) + len(p.MermaidFilesToCreate) + len(p.RedocFilesToCreate))
		diagramCreator(p.FilesToCreate, HttpToFile, progress)
	}
	wg.Wait()
}

// GetFuncMap returns the funcs that are used in diagram generation.
func (p *Generator) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"CreateIntegrationDiagram": p.CreateIntegrationDiagram,
		"CreateSequenceDiagram":    p.CreateSequenceDiagram,
		"CreateParamDataModel":     p.CreateParamDataModel,
		"CreateReturnDataModel":    p.CreateReturnDataModel,
		"CreateTypeDiagram":        p.CreateTypeDiagram,
		"CreateRedoc":              p.CreateRedoc,
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
		"ServiceMetadata":          ServiceMetadata,
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
