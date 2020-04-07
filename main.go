package main

import (
	"net/http"
	"os"

	"github.com/anz-bank/sysl-catalog/pkg/templategeneration"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"gopkg.in/alecthomas/kingpin.v2"
)

//func main() {
//	plantumlService := os.Getenv("SYSL_PLANTUML")
//	if plantumlService == "" {
//		panic("Error: Set SYSL_PLANTUML env variable")
//	}
//	var output string
//	flag.StringVar(&output, "o", "./", "Output directory of documentation")
//	flag.Parse()
//	filename := flag.Arg(0)
//	fs := afero.NewOsFs()
//	m, err := parse.NewParser().Parse(filename, fs)
//	fs = afero.NewMemMapFs()
//	if err != nil {
//		panic(err)
//	}
//	templategeneration.NewProject(filename, "/"+output, plantumlService, logrus.New(), fs, m).ExecuteTemplateAndDiagrams()
//
//	httpFs := afero.NewHttpFs(fs)
//	fileserver := http.FileServer(httpFs.Dir("/"))
//	http.Handle("/", fileserver)
//	http.ListenAndServe(":80", fileserver)
//}

var (
	input = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()

	server     = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	port       = kingpin.Flag("port", "Port to serve on").Short('p').Default(":69").String()
	outputType = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir  = kingpin.Flag("output", "Output directory to generate to").Short('o').String()
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(*input, fs)

	if *server {
		fs = afero.NewMemMapFs()
		*outputDir = "/" + *outputDir
		*outputType = "html"
	}
	if err != nil {
		panic(err)
	}
	println(*outputType)
	templategeneration.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), fs, m).ExecuteTemplateAndDiagrams()
	if *server {
		httpFs := afero.NewHttpFs(fs)
		fileserver := http.FileServer(httpFs.Dir("/"))
		http.Handle("/", fileserver)
		http.ListenAndServe(*port, fileserver)
	}
}
