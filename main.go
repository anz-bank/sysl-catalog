package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/anz-bank/sysl-catalog/pkg/swaggergen"
	"github.com/anz-bank/sysl-catalog/pkg/templategeneration"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/go-openapi/runtime/middleware"
	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	input      = kingpin.Arg("input", "input sysl file to generate documentation for").Required().String()
	server     = kingpin.Flag("serve", "Start a http server and preview documentation").Bool()
	appname    = kingpin.Flag("appname", "Application name for swagger api generation").String()
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

		if *appname != "" {
			// Watch input dir to regenerate swagger
			go watchFile(
				func() {
					m, err := parse.NewParser().Parse(*input, fs)
					if err != nil {
						panic(err)
					}
					// regenerate the swagger for only the first application
					err = swaggergen.GenerateSwagger(*appname, m, httpFileSystem)
					if err != nil {
						log.Printf("failed to generate swagger")
					}
				},
				path.Dir(*input))
			http.HandleFunc("/redoc/", handleRedoc)
			http.HandleFunc("/swagger/", redirectToSwagger)
		}
		http.ListenAndServe(*port, nil)
	} else {
		m, err := parse.NewParser().Parse(*input, fs)
		if err != nil {
			panic(err)
		}
		templategeneration.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), fs, m).ExecuteTemplateAndDiagrams()
	}
}

func handleRedoc(w http.ResponseWriter, r *http.Request) {
	handler := middleware.Redoc(middleware.RedocOpts{Path: "/redoc"}, nil)
	r.URL.Path = strings.TrimRight(r.URL.Path, "/")
	handler.ServeHTTP(w, r)
}

func redirectToSwagger(w http.ResponseWriter, r *http.Request) {
	// Assuming you want to serve a photo at 'images/foo.png'
	fp := "/swagger.json"
	http.Redirect(w, r, fp, http.StatusSeeOther)
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
