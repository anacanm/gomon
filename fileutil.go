package main

import "strings"

// fileShouldBeWatched is a helper function that returns a boolean
// returns whether or not the file extension is in the acceptedFileExtensions string slice
func fileShouldBeWatched(fileName string, acceptedFileExtensions []string) bool {
	for _, v := range acceptedFileExtensions {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}

func filterOutTests(filesWithTests []string) []string {
	// initialize result with a capacity of the length of filesWith Tests
	result := make([]string, 0, len(filesWithTests))

	for _, fileName := range filesWithTests {
		// if fileName does not end with _test.go, add it to the result slice
		if !strings.HasSuffix(fileName, "_test.go") {
			result = append(result, fileName)
		}
	}

	return result
}
