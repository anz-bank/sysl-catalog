package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/diagrams"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sequencediagram"
	"github.com/spf13/afero"
)

type IndexMarkdown struct {
	PackageName string
	App         []*AppMarkdown
}

const IndexMarkdownTemplate = `
| Package |\n| - |
{{range $Package := .}}{{$Package.PackageName}}|
{{end}}
`

type AppMarkdown struct {
	ServiceName string
	Method      string
	Link        string
}

const AppMarkdownTemplate = `
| Service | Method |\n| - |:-:|\n
{{range $App := .}}{{$App.ServiceName}}|({{$App.Method}})[{{$App.Link}}] |{{end}}
`

func main() {
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	var output, packageName string
	flag.StringVar(&output, "o", "./", "Output directory of documentation")
	flag.Parse()
	filename := flag.Arg(0)
	fmt.Println(filename)
	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(filename, fs)
	if err != nil {
		panic(err)
	}
	// README, err := fs.Create(output + "/README.md")
	// if err != nil {
	// 	panic(err)
	// }
	Index := make(map[string]*IndexMarkdown)

	for _, app := range m.Apps {
		if syslutil.HasPattern(app.Attrs, "ignore") {
			continue
		}
		appName := strings.Join(app.Name.GetPart(), "")
		if attr := app.GetAttrs()["package"]; attr != nil {
			packageName = attr.GetS()
		} else {
			packageName = appName
		}
		packageReadmeName := filepath.Join(output, packageName, "README.md")
		MarkdownApp := AppMarkdown{}
		fs.MkdirAll(path.Join(output, packageName), os.ModePerm)
		// if pacakgeREADME, ok = packageReadmes[packageReadmeName]; !ok {
		// 	pacakgeREADME, err = fs.Create(packageReadmeName)
		// 	packageReadmes[packageReadmeName] = pacakgeREADME
		// pacakgeREADME.Write([]byte("| Service | Method |\n| - |:-:|\n"))
		// README.Write([]byte(fmt.Sprintf("[%s](%s) |\n", packageName, packageName)))
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// } else {
		// 	pacakgeREADME = packageReadmes[packageReadmeName]
		// }
		for _, endpoint := range app.Endpoints {
			outputFileName := path.Join(output, packageName, appName+endpoint.Name+".svg")
			MarkdownApp = AppMarkdown{
				ServiceName: appName,
				Method:      endpoint.Name,
				Link:        outputFileName,
			}
			CreateSequenceDiagramFile(
				m,
				fmt.Sprintf("%s <- %s", appName, endpoint.Name),
				outputFileName,
				plantumlService,
				fs)
		}
		if _, ok := Index[packageReadmeName]; !ok {
			Index[packageReadmeName] = &IndexMarkdown{}
			Index[packageReadmeName].App = make([]*AppMarkdown, 10)
		}
		Index[packageReadmeName].App = append(Index[packageReadmeName].App, &MarkdownApp)
	}
	README, err := fs.Create(output + "/README.md")
	if err != nil {
		panic(err)
	}
	IndexTemplate, err := template.New("markdown").Parse(IndexMarkdownTemplate)
	err = IndexTemplate.Execute(README, Index)
	if err != nil {
		panic(err)
	}
	AppTemplate, err := template.New("markdown").Parse(AppMarkdownTemplate)
	if err != nil {
		panic(err)
	}
	for Name, Apps := range Index {
		README, err := fs.Create(Name)
		if err != nil {
			panic(err)
		}
		err = AppTemplate.Execute(README, Apps)
	}
}

// func GenerateMarkdown(index []IndexMarkdown, outputName string, fs){

// }

func CreateSequenceDiagramFile(m *sysl.Module, call, outputFileName, plantumlService string, fs afero.Fs) error {
	sequenceDiagram, err := CreateSequenceDiagram(m, call)
	if err != nil {
		panic(err)
	}
	return diagrams.OutputPlantuml(outputFileName, plantumlService, sequenceDiagram, fs)
}
func CreateSequenceDiagram(m *sysl.Module, call string) (string, error) {
	l := &cmdutils.Labeler{}
	p := &sequencediagram.SequenceDiagParam{}
	p.Endpoints = []string{call}
	p.AppLabeler = l
	p.EndpointLabeler = l
	p.Title = call
	return sequencediagram.GenerateSequenceDiag(m, p, logrus.New())
}
