package main

import "strings"

// fileShouldBeWatched is a helper function that returns a boolean
// returns whether or not the file extension is in the acceptedFileExtensions string slice
// TODO: Add testing
func fileShouldBeWatched(fileName string, acceptedFileExtensions []string) bool {
	for _, v := range acceptedFileExtensions {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}
