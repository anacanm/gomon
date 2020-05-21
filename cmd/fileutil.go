package main

import "strings"

func fileShouldBeWatched(fileName string, acceptedFileExtensions []string) bool {
	for _, v := range acceptedFileExtensions {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}
