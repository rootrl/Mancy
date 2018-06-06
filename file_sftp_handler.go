package main

import (
	"time"
	"log"
	"github.com/pkg/sftp"
)

var (
	writeJobChan  = make(chan bool)
	createJobChan = make(chan bool)
	removeJobChan = make(chan bool)
	renameJobChan = make(chan bool)
	chmodJobChan  = make(chan bool)
	err           error
	sftpClient    *sftp.Client
)

func fileSftpHandler() {
	// Connect
	sftpClient, err = connect("root", "@rootrl13156605816", "117.48.205.33", 22)

	if err != nil {
		log.Fatal("SSH connect error: ", err)
	}

	defer sftpClient.Close()

	// Check if the remote dir exist
	if _, err := sftpClient.Stat(*remoteDir); err != nil {
		panic("Remote dir dose not exist: " + *remoteDir)
	}

	// Create Event handler
	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case createFileName := <-fileCreateEvent:
				log.Print("createFile:" + createFileName)

				if isDir(createFileName) == 1 {
					uploadDirectory(createFileName, *remoteDir)
				} else {
					uploadFile(createFileName, *remoteDir)
				}
			}
		}

		createJobChan <- true
	}()

	// Write event handler
	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case writeFileName := <-fileWriteEvent:
				log.Print("write file: " + writeFileName)

				isDir := isDir(writeFileName)

				if isDir == 1 {
					uploadDirectory(writeFileName, *remoteDir)
				} else if isDir == 0 {
					uploadFile(writeFileName, *remoteDir)
				}
			}
		}

		writeJobChan <- true
	}()

	// Remove event handler
	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case removeFileName := <-fileRemoveEvent:
				log.Print("remove file: " + removeFileName)

				remove(removeFileName, *remoteDir)
			}
		}

		removeJobChan <- true
	}()

	// Rename event handler
	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case renameFileName := <-fileRenameEvent:
				log.Print("rename file: " + renameFileName)
			}
		}

		renameJobChan <- true
	}()

	// Chmod event handler
	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case chmodFileName := <-fileChmodEvent:
				log.Print("chmod file" + chmodFileName)
			}
		}

		chmodJobChan <- true
	}()

	<-writeJobChan
	<-createJobChan
	<-removeJobChan
	<-renameJobChan
	<-chmodJobChan
}
