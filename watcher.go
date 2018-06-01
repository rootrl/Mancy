package main

import (
	"log"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"os"
	"fmt"
)

type Watch struct {
	watch *fsnotify.Watcher;
}

// handler jobs done
var done = make(chan bool)

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
			log.Print("Watching: ", path)
		}

		return nil
	})

	// Handle the watch events
	go eventsHandler(w)

	// Await
	<-done
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

					fmt.Println(ev.Name)

					if err == nil && fi.IsDir() {
						w.watch.Add(ev.Name)
					}
				}

				// write event
				if ev.Op&fsnotify.Write == fsnotify.Write {
				}

				// delete event
				if ev.Op&fsnotify.Remove == fsnotify.Remove {

					fi, err := os.Stat(ev.Name)

					if err == nil && fi.IsDir() {
						w.watch.Remove(ev.Name)
					}
				}

				// Rename
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					w.watch.Remove(ev.Name)
				}
				// Chmod
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				}
			}
		case err := <-w.watch.Errors:
			{
				log.Fatal(err)
				done <- true
				return
			}
		}
	}

	done <- true
}