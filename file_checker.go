package main

import (
	"bufio"
	"io"
	"os"

	"path/filepath"
	"strings"
)

// Check the file required is contains the ignore file
// if contains, return true
func fileChecker(filename string) bool {
	filePathName, err := filepath.Abs("./ignorefile")

	if err != nil {
		return false
	}

	file, err := os.Open(filePathName)

	if err != nil {
		return  false
	}

	defer file.Close()

	fileBuf := bufio.NewReader(file)

	for {
		content, _, endLine := fileBuf.ReadLine()
		if endLine == io.EOF {
			break
		}

		if strings.Contains(filename, string(content)) {
			return  true
		}
	}

	return false
}