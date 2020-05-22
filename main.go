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

	// begin walking at the root of the directory, adding a fsnotify watcher to each file
	// TODO: make file extensions specifiable by a command line flag
	if err := filepath.Walk(".", addWatcherToAppropriateFiles); err != nil {
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

							// if the filepath provided by the user does not have a "/" at the end, as it should to specify that it is a dir
							if !strings.HasSuffix(filePathToLookForGoFiles, "/") {
								filePathToLookForGoFiles += "/"
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
				// TODO: propogate errors rather than just printing them
				fmt.Printf("Error recieved from file watcher: %v", err)
			}
		}
	}()

	<-done

	return nil
}

// addWatcherToAppropriateFiles is run as the WalkFunc for filepath.Walk/2
// addWatcherToAppropriateFiles adds a fsnotify watcher to each file that should be monitored
func addWatcherToAppropriateFiles(path string, info os.FileInfo, err error) error {
	// first, check the error in the function parameters
	if err != nil {
		return err
	}

	// add fsnotify watchers to files (not directories) that should be watched (have specific file extensions)
	if !info.IsDir() {
		// TODO: make file extensions to be watched specified by user
		if fileShouldBeWatched(info.Name(), []string{".go"}) {
			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("Error adding fsnotify watcher to %v: %v", path, err)
			}

		}
	}
	return nil
}
