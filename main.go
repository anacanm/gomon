package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if err := WatchFiles(); err != nil {
		log.Fatalf("Error HERE: %v", err)
	}
}

var watcher *fsnotify.Watcher

// WatchFiles watches the appropriate files, sending events or errors as they occur on files
func WatchFiles() error {

	// creates a new file watcher
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("Error creating new fsnotify watcher: %v", err)
	}
	defer watcher.Close()

	// get the extensionsToBeWatched specified by the user using -flags
	extensionsToBeWatched := fileExtensionsToBeWatched()

	// begin walking at the root of the directory, adding a fsnotify watcher to each file
	if err := filepath.Walk(".", getWatcherWalkFunc(extensionsToBeWatched)); err != nil {
		return fmt.Errorf("Error walking filepath: %v", err)
	}

	done := make(chan bool)

	// start a new goroutine to listen for events
	go func() {
		for {
			select {
			// listen for a file event or an error
			case event := <-watcher.Events:
				if event.Op == fsnotify.Write {
					// if a write event is received, then a file that we added a watcher to was modified
					// therefore, I should restart the go project by running the go files in the specified directory
					go func() {
						// goFileMatches is a slice of strings that are the names of the files in the specified directory that are go files that need to be run
						var filePathToLookForGoFiles string
						if len(os.Args) > 1 {
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

						goFileMatches, err := filepath.Glob(filePathToLookForGoFiles + "*.go")
						if err != nil {
							fmt.Printf("Error getting filematches: %v", err)
						}

						// add "run" to the begginning of the filenames, to be used below in the go run command
						goFileMatches = append([]string{"run"}, goFileMatches...)

						// create the command to be run: go run FILE_LOCATIONS
						// 		ex: go run cmd/main.go cmd/fileutil.go
						cmd := exec.Command("go", goFileMatches...)

						// run the command, and access the returned combined standard output and standard error
						stdoutStderr, err := cmd.CombinedOutput()
						if err != nil {
							fmt.Printf("GOMON REPORTED ERROR: %v\n\n", err)
						}

						// redirects stdOut and stdErr of file to stdOut & stdErr of gomon
						fmt.Printf("%s", stdoutStderr)
					}()
				}
			case err := <-watcher.Errors:
				fmt.Printf("Error recieved from file watcher: %v", err)
			}
		}
	}()

	<-done

	return nil
}

func getWatcherWalkFunc(extensionsToBeWatched []string) filepath.WalkFunc {
	// since filepath.Walk requires a filepath.Walkfunc function, the callback function that I provide needs to have that specific signature, ie. I can't pass another param
	// therefore, getWatcherWalkFunc needs to return a closure that has access to extensionsToBeWatched, since the flags can (and should) only  be parsed one time

	// the below anonymous function is run as the WalkFunc for filepath.Walk/2
	// the below anonymous function adds a fsnotify watcher to each file that should be monitored according to extensionsToBeWatched
	return func(path string, info os.FileInfo, err error) error {
		// first, check the error in the function parameters
		if err != nil {
			return err
		}

		// add fsnotify watchers to files (not directories) that should be watched (have been specified by the user with -flags)
		if !info.IsDir() {
			if fileShouldBeWatched(info.Name(), extensionsToBeWatched) {
				if err := watcher.Add(path); err != nil {
					return fmt.Errorf("Error adding fsnotify watcher to %v: %v", path, err)
				}

			}
		}
		return nil
	}
}
