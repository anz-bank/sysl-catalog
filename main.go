package main

import (
	"fmt"
	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/radovskyb/watcher"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input      = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	server     = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	port       = kingpin.Flag("port", "Port to serve on").Short('p').Default(":6900").String()
	outputType = kingpin.Flag("type", "Type of output").HintOptions("html", "markdown").Default("markdown").String()
	outputDir  = kingpin.Flag("output", "Output directory to generate to").Short('o').String()
	verbose     = kingpin.Flag("verbose", "Verbose logs").Short('v').Bool()
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}
	fs := afero.NewOsFs()
	log := logrus.New()
	if *verbose{
		log.SetLevel(logrus.InfoLevel)
	}else{
		log.SetLevel(logrus.ErrorLevel)
	}
	m, err := parse.NewParser().Parse(*input, fs)
	if err != nil{
		log.Fatal(err)
	}
	if *server {
		httpfilesystem := afero.NewMemMapFs()
		httpFs := afero.NewHttpFs(httpfilesystem)
		fileserver := http.FileServer(httpFs.Dir("/"))
		go watchFile(func(){
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Error:", r)
				}
			}()
			m, err := parse.NewParser().Parse(*input, fs)
			if err != nil{
				log.Fatal(err)
			}
			catalog.NewProject(*input, "/" + *outputDir, plantumlService, "html", log, m).SetOutputFs(httpfilesystem).ExecuteTemplateAndDiagrams()
			fmt.Println("Done Regenerating")
		}, path.Dir(*input))
		fmt.Println("Serving on http://localhost"+*port)
		log.Fatal(http.ListenAndServe(*port, fileserver))
		select{}
	}
	project := catalog.NewProject(*input, *outputDir, plantumlService, *outputType, log, m)
	project.SetOutputFs(fs).ExecuteTemplateAndDiagrams()
}

func watchFile(action func(), files ...string) {
	w := watcher.New()
	// Only notify rename and move events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event)
				action()
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