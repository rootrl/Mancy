package main

import (
	"os"
	"strings"
	"path/filepath"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// Check the filePath if is a dir
func isDir(fileName string) int {
	fileInfo, err := os.Stat(fileName)

	if err != nil {
		log.Println("file dose not exist: " + fileName)
		return -1
	}

	if fileInfo.IsDir() {
		return 1
	} else {
		return 0
	}
}

// Get the change part of filePath
func getChangedFileName(fileName string) string {
	fileName = strings.Replace(fileName, localDir, "", -1)
	return filepath.ToSlash(fileName)
}

// Check and generate config file
func checkConfigFile() {

	// Check
	if _, err := os.Stat(ConfigFilePathName); !os.IsNotExist(err) {
		return
	}

	// Generate if config not exist

	// config template
	var tpl string = `{
	"localDir": "./",
	"remoteDir": "Remote dir (absolute path like: /root/a)",

	"sshHost": "Your host",
	"sshPort": 22,
	"sshUserName": "Your user name",
	"sshPassword": "Your password",

	"ignoreFiles": [
		".git",
		".idea",
		".swp",
		".swx",
		"___jb_old___",
		"___jb_tmp___"
	]
}`

	if err := ioutil.WriteFile(ConfigFilePathName, []byte(tpl), 0755); err != nil {
		panic(err)
	}

	fmt.Println("\n", `Your are first time run this script, The config file 'config,json' not found, And already generated for you.`,
		`Please reset the config item as your own infomation suck as ssh host & account, then run script again.`, "\n")
	os.Exit(0)
}

// Parse config file
func parseConfigFile() {

	configJson, err := ioutil.ReadFile(ConfigFilePathName)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configJson, &config)

	if err != nil {
		panic(err)
	}
}
