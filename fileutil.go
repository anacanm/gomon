package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

// filterOutTests returns a NEW string slice of go filenames, with test files filtered out
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

// getFilesToRun returns a string slice of go filenames in the specified (default: root of where gomon is being run) directory that should be run
func getFilesToRun() []string {
	// goFileMatches is a slice of strings that are the names of the files in the specified directory that are go files that need to be run
	var filePathToLookForGoFiles string
	if len(os.Args) > 1 {
		// the user has the option to specify where go run should be ran (ex, cmd)
		// if they did provide at least one command line arg, then we need to check to see if a file location was provided
		filePathToLookForGoFiles = os.Args[1]

		// if the filepath provided by the user does not have a "/" at the end (as it should) to specify that it is a dir
		// then add the "/"
		if !strings.HasSuffix(filePathToLookForGoFiles, "/") {
			filePathToLookForGoFiles += "/"
		}

		// if the first command line argument begins with a "-", then it is a flag, and gomon should be run in the root dir
		if strings.HasPrefix(filePathToLookForGoFiles, "-") {
			filePathToLookForGoFiles = ""
		}

	}
	// find all files that have the .go file extensions where the user specified to look
	// we need to run them like "go run main.go flagutil.go fileutil.go", creating our own wildcard (*) functionality
	goFileMatches, err := filepath.Glob(filePathToLookForGoFiles + "*.go")
	if err != nil {
		fmt.Printf("Error getting filematches: %v", err)
	}

	goFileMatchesWithoutTests := filterOutTests(goFileMatches)

	return goFileMatchesWithoutTests
}
