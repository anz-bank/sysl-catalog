package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
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
	var fs afero.Fs
	fs = afero.NewOsFs()
	m, err := parse.NewParser().Parse(*input, fs)
	if *server {
		watchFile(func(){
			*outputDir = "/" + *outputDir
			*outputType = "html"
			handler := catalog.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), m).HTTPHandler()
			go http.ListenAndServe(*port, handler)
		}, path.Dir(*input))

	} else {
		if err != nil {
			panic(err)
		}
		project := catalog.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), m)
		project.SetOutputFs(fs).ExecuteTemplateAndDiagrams()
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
