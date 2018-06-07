package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"time"
	"path/filepath"
)

// Package const
// Config file path name
const (
	ConfigFilePathName = "./mancy_config.json"
)

// Package variables
var (
	// Local dir name witch will be watched
	localDir string

	// Local dir name witch will be watched
	remoteDir string

	// Config data struct
	config Config

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

// Config file struct
type Config struct{
	// filePath
	LocalDir string `localDir`
	RemoteDir string `remoteDir`

	// ssh
	SshHost string `sshHost`
	SshPort int `sshPort`
	SshUserName string `sshUserName`
	SshPassword string `sshPassword`

	// ignoreFiles
	IgnoreFiles []string `ignoreFiles`
}

// Init
func init() {
	// Reset log format
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Check and generate config file
	checkConfigFile()

	// Parse config file
	parseConfigFile()

	// Set localDir as ABS path
	localDir, err = filepath.Abs(config.LocalDir)

	// Set remote dir
	remoteDir = config.RemoteDir

	if err != nil {
		log.Fatal("Init localFilePath error: ")
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
		w.watchDir(localDir)
		watcherHandlerDone <- true
	}()

	// handle the file events
	go func() {
		// Handle file with sftp (autoUpload changes)
		// And you can change the handler whatever you need like rsync
		fileSftpHandler()

		fileHandlerDone <- true
	}()

	// Waiting job done
	<-watcherHandlerDone
	<-fileHandlerDone
}
