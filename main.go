package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/anz-bank/sysl-catalog/pkg/templategeneration"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input      = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	server     = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	port       = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir  = kingpin.Flag("output", "Output directory to generate to").Short('o').String()
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		panic("Error: Set SYSL_PLANTUML env variable")
	}
	var fs, httpFileSystem afero.Fs
	fs = afero.NewOsFs()
	httpFileSystem = afero.NewMemMapFs()
	if *server {
		*outputDir = "/" + *outputDir
		*outputType = "html"
		// Watch our input dir for changes and
		go watchFile(
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered in f", r)
					}
				}()
				m, err := parse.NewParser().Parse(*input, fs)
				if err != nil {
					panic(err)
				}
				templategeneration.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), httpFileSystem, m).SetServerMode().ExecuteTemplateAndDiagrams()
			},
			path.Dir(*input))
		httpFs := afero.NewHttpFs(httpFileSystem)
		fileserver := http.FileServer(httpFs.Dir("/"))
		http.Handle("/", fileserver)
		http.ListenAndServe(*port, fileserver)
	} else {
		m, err := parse.NewParser().Parse(*input, fs)
		if err != nil {
			panic(err)
		}
		templategeneration.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), fs, m).ExecuteTemplateAndDiagrams()
	}
}

func watchFile(action func(), files ...string) {
	w := watcher.New()
	// Only notify rename and move events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.Event:
				action()
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()
	//// Watch test_folder recursively for changes.
	for _, file := range files {
		if err := w.AddRecursive(file); err != nil {
			log.Fatalln(err)
		}
	}
	go func() {
		w.Wait()
		w.TriggerEvent(watcher.Write, nil)
	}()
	//// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
