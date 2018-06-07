package main

import (
	"strings"
)

// Check if the file is contains the ignore file
// Return true if contains
func fileChecker(filename string) bool {

	for _, ignoreFile := range config.IgnoreFiles {
		if strings.Contains(filename, string(ignoreFile)) {
			return true
		}
	}

	return false

}