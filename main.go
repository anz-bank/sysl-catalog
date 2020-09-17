package main

import (
	"fmt"
	"io/ioutil"
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
	watch "github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input             = kingpin.Arg("input", "Input sysl file to generate documentation for").Required().String()
	plantUMLoption    = kingpin.Flag("plantuml", "Plantuml service to use").String()
	port              = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType        = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir         = kingpin.Flag("output", "OutputDir directory to generate to").Short('o').String()
	verbose           = kingpin.Flag("verbose", "Verbose logs").Short('v').Bool()
	templates         = kingpin.Flag("templates", "Custom templates to use, separated by a comma").String()
	outputFileName    = kingpin.Flag("outputFileName", "Output file name for pages; {{.Title}}").Default("").String()
	server            = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	noCSS             = kingpin.Flag("noCSS", "Disable adding css to served html").Bool()
	disableLiveReload = kingpin.Flag("disableLiveReload", "Disable live reload").Default("false").Bool()
	noImages          = kingpin.Flag("noImages", "Disable image generation").Default("false").Bool()
	embed             = kingpin.Flag("embed", "Embed images instead of creating svgs").Default("false").Bool()
	enableMermaid     = kingpin.Flag("mermaid", "Enable mermaid for diagram generation instead of Plantuml (Not currently supported)").Default("false").Bool()
	enableRedoc       = kingpin.Flag("redoc", "Enable Redoc generation for specs imported from OpenAPI. Must be run on a git repo.").Default("false").Bool()
	imageDest         = kingpin.Flag("imageDest", "Optional image directory destination (can be outside of the path provided using --output)").String()
)

func main() {
	kingpin.Parse()

	logger := setupLogger()
	plantUMLService := plantUMLService()
	fs := afero.NewOsFs()

	if !*server {
		m, err := parseSyslFile(".", *input, fs, logger)
		if err != nil {
			logger.Fatal(err)
		}

		catalog.NewProject(*input, plantUMLService, *outputType, logger, m, fs, *outputDir).
			SetOptions(*noCSS, *noImages, *embed, *enableRedoc, *enableMermaid, *outputFileName, *imageDest).
			WithTemplateFs(fs, strings.Split(*templates, ",")...).
			Run()

		return
	}

	if *outputType == "markdown" {
		logger.Warn("Server mode uses html as output type by default")
	}
	if *outputDir != "" {
		logger.Warn("OutputDir is ignored in server mode")
	}

	handler := catalog.NewProject(*input, plantUMLService, "html", logger, nil, nil, "").
		SetOptions(*noCSS, *noImages, *embed, *enableRedoc, *enableMermaid, *outputFileName, *imageDest).
		WithTemplateFs(fs, strings.Split(*templates, ",")...).
		ServerSettings(*noCSS, !*disableLiveReload, true)

	// Generate initial files so we don't get 404
	handler.Run()

	logrus.SetOutput(ioutil.Discard)

	go watcher.WatchFile(func(i interface{}) {
		logger.Info("Regenerating...")
		m, err := func() (m *sysl.Module, err error) {
			defer func() {
				if r := recover(); r != nil {
					m = nil
					err = fmt.Errorf("%s", r)
				}
			}()
			if handler.RootModule == nil {
				m, err = parseSyslFile(".", *input, fs, logger)
			} else {
				var changedModule *sysl.Module
				relativeChangedFilePath := "." + strings.TrimPrefix(i.(watch.Event).Path, "/usr")
				changedModule, err = parseSyslFile(".", relativeChangedFilePath, fs, logger)
				if err == nil {
					m = overwriteSyslModules(handler.RootModule, changedModule)
				}
			}

			return
		}()
		if err != nil {
			fmt.Println(err)
		}
		handler.Update(m, err)
		livereload.ForceRefresh()
		logger.Info(i)
		logger.Info("Done Regenerating")
	}, path.Dir(*input))

	http.Handle("/", handler)
	livereload.Initialize()
	http.HandleFunc("/livereload.js", livereload.ServeJS)
	http.HandleFunc("/livereload", livereload.Handler)
	fmt.Println("Serving on http://localhost" + *port)
	logger.Fatal(http.ListenAndServe(*port, nil))
}

// overwriteSyslModules takes two sysl modules and overwrites one with the other
// TODO: Handle app definitions from multiple files
func overwriteSyslModules(existing *sysl.Module, overwrite *sysl.Module) *sysl.Module {
	for k, v := range overwrite.Apps {
		existing.Apps[k] = v
	}
	return existing
}

func plantUMLService() string {
	plantUMLService := os.Getenv("SYSL_PLANTUML")
	if *plantUMLoption != "" {
		plantUMLService = *plantUMLoption
	}
	if plantUMLService == "" && !*enableMermaid {
		log.Fatal("Error: Set SYSL_PLANTUML env variable or --plantuml flag")
	}
	return plantUMLService
}

func setupLogger() *logrus.Logger {
	logger := logrus.New()
	if *verbose {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.ErrorLevel)
		logrus.SetLevel(logrus.ErrorLevel)
	}
	return logger
}

func parseSyslFile(root string, filename string, fs afero.Fs, logger *logrus.Logger) (*sysl.Module, error) {
	logger.Info("Parsing...")
	start := time.Now()
	m, _, err := loader.LoadSyslModule(root, filename, fs, logger)
	elapsed := time.Since(start)
	logger.Info("Done, time elapsed: ", elapsed)
	return m, err
}
