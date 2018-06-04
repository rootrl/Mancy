package main

import (
	"log"
	"os"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
)

type Watch struct {
	watch *fsnotify.Watcher;
}

// handler jobs done
var eventDone = make(chan bool)

// Watch a directory
func (w *Watch) watchDir(dir string) {

	// Walk all directory
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		// Just watch directory(all child can be watched)
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				log.Fatal(err)
			}
			err = w.watch.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	log.Print("Watching: ", dir)

	// Handle the watch events
	go eventsHandler(w)

	// Await
	<-eventDone
}

// Handle the watch events
func eventsHandler(w *Watch) {
	for {
		select {
		case ev := <-w.watch.Events:
			{
				// Create event
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(ev.Name)

					if !fileChecker(ev.Name) {
						fileCreateEvent <- ev.Name
					}

					if err == nil && fi.IsDir() {
						w.watch.Add(ev.Name)
					}
				}

				// write event
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if !fileChecker(ev.Name) {
						fileWriteEvent <- ev.Name
					}
				}

				// delete event
				if ev.Op&fsnotify.Remove == fsnotify.Remove {

					fi, err := os.Stat(ev.Name)

					if err == nil && fi.IsDir() {
						w.watch.Remove(ev.Name)
					}

					if !fileChecker(ev.Name) {
						fileRemoveEvent <- ev.Name
					}
 				}

				// Rename
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					w.watch.Remove(ev.Name)

					if !fileChecker(ev.Name) {
						fileRenameEvent <- ev.Name
					}
				}
				// Chmod
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if !fileChecker(ev.Name) {
						fileChmodEvent <- ev.Name
					}
				}
			}
		case err := <-w.watch.Errors:
			{
				log.Fatal(err)
				eventDone <- true
				return
			}
		}
	}

	eventDone <- true
}
