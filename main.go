package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"
	"github.com/gohugoio/hugo/livereload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input             = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	plantUMLoption    = kingpin.Flag("plantuml", "plantuml service to use").String()
	port              = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType        = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir         = kingpin.Flag("output", "OutputDir directory to generate to").Short('o').String()
	verbose           = kingpin.Flag("verbose", "Verbose logs").Short('v').Bool()
	templates         = kingpin.Flag("templates", "custom templates to use, separated by a comma").String()
	outputFileName    = kingpin.Flag("outputFileName", "output file name for pages; {{.Title}}").Default("").String()
	server            = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	noCSS             = kingpin.Flag("noCSS", "disable adding css to served html").Bool()
	disableLiveReload = kingpin.Flag("disableLiveReload", "diable live reload").Default("false").Bool()
	noImages          = kingpin.Flag("noImages", "don't create images").Default("false").Bool()
	embed             = kingpin.Flag("embed", "Embed images instead of creating svgs").Default("false").Bool()
	enableMermaid     = kingpin.Flag("mermaid", "use mermaid diagrams where possible").Default("false").Bool()
	enableRedoc       = kingpin.Flag("redoc", "generate redoc for specs imported from openapi. Must be run on a git repo.").Default("false").Bool()
	imageDest         = kingpin.Flag("imageDest", "Optional image directory destination (can be outside output)").String()
)

func main() {
	kingpin.Parse()

	plantumlService := os.Getenv("SYSL_PLANTUML")
	if *plantUMLoption != "" {
		plantumlService = *plantUMLoption
	}
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable or --plantuml flag")
	}

	fs := afero.NewOsFs()

	log := logrus.New()
	if *verbose {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.WarnLevel)
		logrus.SetLevel(logrus.WarnLevel)
	}

	if !*server {
		log.Info("Parsing")
		start := time.Now()
		m, _, err := loader.LoadSyslModule(".", *input, fs, log)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		log.Info("Done, time elapsed: ", elapsed)

		catalog.NewProject(*input, plantumlService, *outputType, log, m, fs, *outputDir, *enableMermaid).
			SetOptions(*noCSS, *noImages, *embed, *enableRedoc, *outputFileName, *imageDest).
			WithTemplateFs(fs, strings.Split(*templates, ",")...).
			Run()
		return
	}

	if *outputType == "markdown" {
		log.Warn("Server mode uses html as output type by default")
	}

	handler := catalog.
		NewProject(*input, plantumlService, "html", log, nil, nil, "", *enableMermaid).
		SetOptions(*noCSS, *noImages, *embed, *enableRedoc, *outputFileName, *imageDest).
		WithTemplateFs(fs, strings.Split(*templates, ",")...).
		ServerSettings(*noCSS, !*disableLiveReload, true)
	fmt.Println("Serving on http://localhost" + *port)

	logrus.SetOutput(ioutil.Discard)

	go watcher.WatchFile(func(i interface{}) {
		log.Info("Regenerating")
		m, err := func() (m *sysl.Module, err error) {
			defer func() {
				if r := recover(); r != nil {
					m = nil
					err = fmt.Errorf("%s", r)
				}
			}()
			log.Info("Parsing")
			m, _, err = loader.LoadSyslModule("", *input, fs, log)
			log.Info("Done Parsing")
			return
		}()
		if err != nil {
			fmt.Println(err)
		}
		handler.Update(m, err)
		livereload.ForceRefresh()
		log.Info(i)
		log.Info("Done Regenerating")
	}, path.Dir(*input))

	livereload.Initialize()
	http.HandleFunc("/livereload.js", livereload.ServeJS)
	http.HandleFunc("/livereload", livereload.Handler)
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(*port, nil))
	select {}
}
