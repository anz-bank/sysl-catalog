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
	"os/exec"
	"path"
	"runtime"
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
)

func main() {
	kingpin.Parse()
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}
	fs := afero.NewOsFs()
	m, err := parse.NewParser().Parse(*input, fs)
	if err != nil{
		log.Fatal(err)
	}
	if *server {
		go watchFile(func(){
			http.ListenAndServe(*port, catalog.NewProject(*input, "/" + *outputDir, plantumlService, "markdown", logrus.New(), m))
		}, path.Dir(*input))
		time.Sleep(2 * time.Second)
		openBrowser("http://localhost" + *port)
		select{}
	}
	project := catalog.NewProject(*input, *outputDir, plantumlService, *outputType, logrus.New(), m)
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
				go action()
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

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}