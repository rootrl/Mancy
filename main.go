package main;

import (
	"github.com/fsnotify/fsnotify"
	"flag"
	"log"
)

// Package variables
var (
	// Local dir name witch will be watched
	localDir = flag.String("localDir", "./", "Set the local directory witch will be watched")

	// Local dir name witch will be watched
	remoteDir = flag.String("remoteDir", "./", "Set the remote directory that accept the changes")
)

// Init
func init(){
	// Reset log format
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	watch, _ := fsnotify.NewWatcher()

	w := Watch{
		watch: watch,
	}

	// Watch the local directory
	w.watchDir(*localDir)
}