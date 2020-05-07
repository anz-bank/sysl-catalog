package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/gohugoio/hugo/livereload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input      = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	server     = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	port       = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir  = kingpin.Flag("output", "OutputDir directory to generate to").Short('o').String()
	verbose    = kingpin.Flag("verbose", "Verbose logs").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
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
	m, err := parse.NewParser().Parse(*input, fs)
	if err != nil {
		log.Fatal(err)
	}
	if *server {
		handler := catalog.NewProject(*input, plantumlService, "html", log, m, nil, "").
			ServerSettings(true, true, true)

		go watcher.WatchFile(func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Error:", r)
				}
			}()
			m, err := parse.NewParser().Parse(*input, fs)
			if err != nil {
				log.Fatal(err)
			}
			handler.Update(m)
			livereload.ForceRefresh()
			fmt.Println("Done Regenerating")
		}, path.Dir(*input))
		fmt.Println("Serving on http://localhost" + *port)

		livereload.Initialize()
		http.HandleFunc("/livereload.js", livereload.ServeJS)
		http.HandleFunc("/livereload", livereload.Handler)
		http.Handle("/", handler)
		log.Fatal(http.ListenAndServe(*port, nil))
		select {}
	}
	catalog.NewProject(*input, plantumlService, *outputType, log, m, fs, *outputDir).Run()
}
