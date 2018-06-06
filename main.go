package main;

import (
	"github.com/fsnotify/fsnotify"
	"flag"
	"log"
	"time"
	"path/filepath"
)

// Package variables
var (
	// Local dir name witch will be watched
	localDir = flag.String("local-dir", "./", "Set the local directory witch will be watched")

	// Local dir name witch will be watched
	remoteDir = flag.String("remote-dir", "/root/a", "Set the remote directory that accept the changes")

	// LocalDir with slash separator
	slashLocalDir string

	// Global chan variables
	// file_watcher will write the chan and file_handle will read the chan
	// create file
	fileCreateEvent = make(chan string)

	// write
	fileWriteEvent = make(chan string)

	// remove
	fileRemoveEvent = make(chan string)

	// rename
	fileRenameEvent = make(chan string)

	// chmod
	fileChmodEvent = make(chan string)

	// watchMainJob chan
	watcherHandlerDone = make(chan bool)

	// fileHandleMainJob chan
	fileHandlerDone = make(chan bool)

	// timeout for watcher event
	fileHandleTimeOut = time.Second * 4
)

// Init
func init() {
	// Reset log format
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Set localDir as ABS path
	*localDir, err = filepath.Abs(*localDir)

	// slashLocalDir := filepath.ToSlash(*localDir)

	if err != nil {
		log.Fatal("Init localFilePath error: " )
		panic(err)
	}
}

func main() {
	watch, _ := fsnotify.NewWatcher()

	w := Watch{
		watch: watch,
	}

	// Watch the local directory
	go func() {
		w.watchDir(*localDir)
		watcherHandlerDone <- true
	}()

	// handle the file events
	go func() {
		// Handle file with sftp (autoLoad changes)
		fileSftpHandler()
		fileHandlerDone <- true
	}()

	// Waiting job done
	<-watcherHandlerDone
	<-fileHandlerDone
}