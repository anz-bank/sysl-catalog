package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/buger/goterm"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/gohugoio/hugo/livereload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input             = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	port              = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType        = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir         = kingpin.Flag("output", "OutputDir directory to generate to").Short('o').String()
	verbose           = kingpin.Flag("verbose", "Verbose logs").Short('v').Bool()
	projectTemplate   = kingpin.Flag("projectTemplate", "projectTemplate filname to use").String()
	packageTemplate   = kingpin.Flag("packageTemplate", "packageTemplate filname to use").String()
	server            = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	noCSS             = kingpin.Flag("noCSS", "disable adding css to served html").Bool()
	disableLiveReload = kingpin.Flag("disableLiveReload", "diable live reload").Default("false").Bool()
	noImages          = kingpin.Flag("noImages", "don't create images").Default("false").Bool()
	enableMermaid     = kingpin.Flag("mermaid", "use mermaid diagrams where possible").Default("false").Bool()
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
	var f1, f2 afero.File
	var err error
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}
	fs := afero.NewOsFs()
	log := logrus.New()
	if *verbose {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.ErrorLevel)
		logrus.SetLevel(logrus.ErrorLevel)
	}
	if *projectTemplate != "" {
		f1, err = fs.Open(*projectTemplate)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *packageTemplate != "" {
		f2, err = fs.Open(*packageTemplate)
		if err != nil {
			log.Fatal(err)
		}
	}
	if !*server {
		m, err := parse.NewParser().Parse(*input, fs)
		if err != nil {
			log.Fatal(err)
		}
		catalog.NewProject(*input, plantumlService, *outputType, log, m, fs, *outputDir, *enableMermaid).
			WithTemplateFiles(f1, f2).
			SetOptions(*noCSS, *noImages).
			Run()
		return
	}

	handler := catalog.
		NewProject(*input, plantumlService, "html", log, nil, nil, "", *enableMermaid).
		WithTemplateFiles(f1, f2).
		ServerSettings(*noCSS, !*disableLiveReload, true)
	goterm.Clear()
	PrintToPosition(1, "Serving on http://localhost"+*port)
	logrus.SetOutput(ioutil.Discard)
	go watcher.WatchFile(func(i interface{}) {
		PrintToPosition(3, "Regenerating")
		m, err := func() (m *sysl.Module, err error) {
			defer func() {
				if r := recover(); r != nil {
					m = nil
					err = fmt.Errorf("%s", r)
				}
			}()
			m, err = parse.NewParser().Parse(*input, fs)
			return
		}()
		if err != nil {
			PrintToPosition(4, err)
		}
		handler.Update(m, err)
		livereload.ForceRefresh()
		PrintToPosition(2, i)
		PrintToPosition(4, goterm.RESET_LINE)
		PrintToPosition(3, goterm.RESET_LINE)
		PrintToPosition(3, "Done Regenerating")
	}, path.Dir(*input))
	livereload.Initialize()
	http.HandleFunc("/livereload.js", livereload.ServeJS)
	http.HandleFunc("/livereload", livereload.Handler)
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(*port, nil))
	select {}
}

func PrintToPosition(y int, i interface{}) {
	goterm.MoveCursor(1, y)
	goterm.Print(i)
	goterm.Flush()
}
