package watcher

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"time"
)

func WatchFile(action func(), files ...string) {
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
