package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"

	"github.com/anz-bank/sysl/pkg/mod"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/gohugoio/hugo/livereload"
	watch "github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	runCmd            = kingpin.Command("run", "Run the generator")
	input             = runCmd.Arg("input", "Input sysl file to generate documentation for").Required().String()
	plantUMLoption    = runCmd.Flag("plantuml", "Plantuml service to use").String()
	port              = runCmd.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType        = runCmd.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir         = runCmd.Flag("output", "OutputDir directory to generate to").Short('o').String()
	verbose           = runCmd.Flag("verbose", "Verbose logs").Short('v').Bool()
	templates         = runCmd.Flag("templates", "Custom templates to use, separated by a comma, or 'mermaid' or 'plantuml' for defaults").String()
	outputFileName    = runCmd.Flag("outputFileName", "Output file name for pages; {{.Title}}").Default("").String()
	server            = runCmd.Flag("serve", "Start a http server and preview documentation").Bool()
	noCSS             = runCmd.Flag("noCSS", "Disable adding css to served html").Bool()
	disableLiveReload = runCmd.Flag("disableLiveReload", "Disable live reload").Default("false").Bool()
	modCmd            = kingpin.Command("mod", "sysl modules")
	cmd               = modCmd.Arg("cmd", "get or update").String()
	repo              = modCmd.Arg("repo", "repo to get").String()
)

func main() {
	kingpin.Parse()

	logger := setupLogger()
	plantUMLService := plantUMLService()
	fs := afero.NewOsFs()
	retr, err := mod.Retriever(afero.NewOsFs())
	if err != nil {
		logger.Fatal(err)
	}
	if *cmd != "" {
		if err := retr.Command(*cmd, *repo); err != nil {
			logrus.Error(err)
		}
		return
	}
	if !*server {
		m, err := parseSyslFile(".", *input, fs, logger)
		if err != nil {
			logger.Fatal(err)
		}

		catalog.NewProject(*input, plantUMLService, *outputType, logger, m, fs, *outputDir).
			SetOptions(*noCSS, *outputFileName, "/").
			WithRetriever(retr).
			AutomaticTemplates(fs, strings.Split(*templates, ",")...).
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
		SetOptions(*noCSS, *outputFileName, "").
		WithRetriever(retr).
		AutomaticTemplates(fs, strings.Split(*templates, ",")...).
		ServerSettings(*noCSS, !*disableLiveReload, true)

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
				if err != nil {
					return nil, err
				}
			} else {
				var changedModule *sysl.Module
				wd, _ := os.Getwd()
				relativeChangedFilePath := "." + strings.TrimPrefix(i.(watch.Event).Path, wd)
				changedModule, err = parseSyslFile(".", relativeChangedFilePath, fs, logger)
				if err == nil {
					m = overwriteSyslModules(handler.RootModule, changedModule)
				} else {
					return nil, err
				}
			}
			return
		}()
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
	retr, err := mod.Retriever(fs)
	if err != nil {
		return nil, err
	}
	m, err := parse.NewParser().Parse(path.Join(root, filename), retr)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	logger.Info("Done, time elapsed: ", elapsed)
	return m, err
}
