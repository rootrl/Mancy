package main

import (
	"time"
	"log"
)

var (
	writeJobChan = make(chan bool)
	createJobChan = make(chan bool)
	removeJobChan = make(chan bool)
	renameJobChan = make(chan bool)
	chmodJobChan = make(chan bool)
)

func fileHandler() {

	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case createFileName := <-fileCreateEvent:
				log.Print("createFile:" + createFileName)
			}
		}

		createJobChan <- true
	}()

	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case writeFileName := <-fileWriteEvent:
				log.Print("write file: " + writeFileName)
			}
		}

		writeJobChan <- true
	}()

	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case removeFileName := <-fileRemoveEvent:
				log.Print("remove file" + removeFileName)
			}
		}

		removeJobChan <- true
	}()

	go func() {
		for {
			select {

			case <-time.After(fileHandleTimeOut):
				//timeout for 3 seconds

			case renameFileName := <-fileRenameEvent:
				log.Print("rename file" + renameFileName)
			}
		}

		renameJobChan <- true
	}()

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
