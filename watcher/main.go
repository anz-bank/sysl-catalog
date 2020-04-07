package main

import (
	"fmt"
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func watchFile(files ...string) {
	w := watcher.New()
	// Only notify rename and move events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Write)
	go func() {
		for {
			select {
			case event := <-w.Event:
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

	//// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
