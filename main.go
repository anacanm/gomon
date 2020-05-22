package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// first, get a slice of filenames that I should run, according to the user
	filesToRun := getFilesToRun()

	// start running the files in a separate goroutine

	fmt.Printf("starting gomon... running %v\n\n", filesToRun)
	go startFiles(filesToRun)

	if err := WatchFiles(filesToRun); err != nil {
		log.Fatalf("Error HERE: %v", err)
	}
}

var watcher *fsnotify.Watcher

// WatchFiles watches the appropriate files, sending events or errors as they occur on files
func WatchFiles(filesToRun []string) error {

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
				if isModEvent(event.Op) {
					// if a mod event is received, then a file that I added a watcher to was modified
					// therefore, I should restart the go project by running the go files in the specified directory
					go startFiles(filesToRun)
				}
			case err := <-watcher.Errors:
				fmt.Printf("Error recieved from file watcher: %v", err)
			}
		}
	}()

	<-done

	return nil
}

// getWatcherWalkFunc returns an anonymous filepath.WalkFunc that has access to extensionsToBeWatched
// this closure functionality allows the inner function to access the extensionsToBeWatched while still keeping the necessary signature
func getWatcherWalkFunc(extensionsToBeWatched []string) filepath.WalkFunc {
	// since filepath.Walk requires a filepath.Walkfunc function, the callback function that I provide needs to have that specific signature, ie. I can't pass another param
	// therefore, getWatcherWalkFunc needs to return a closure that has access to extensionsToBeWatched, since the flags can (and should) only  be parsed one time

	// the below anonymous function is run as the WalkFunc for filepath.Walk
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

// startFiles starts the given files, connects the stdOut and stdErr to current stdOut and stdErr
func startFiles(filesToRun []string) {
	// add "run" to the begginning of the filenames, to be used below in the go run command
	runFileCommand := append([]string{"run"}, filesToRun...)

	// create the command to be run: go run FILE_LOCATIONS
	// 		ex: go run cmd/main.go cmd/fileutil.go
	cmd := exec.Command("go", runFileCommand...)

	// run the command, and access the returned combined standard output and standard error
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("GOMON REPORTED ERROR: %v\n\n", err)
	}

	// redirects stdOut and stdErr of file to stdOut & stdErr of gomon
	fmt.Printf("%s", stdoutStderr)
}

// isModEvent returns true if the incoming event should result in restarting the files according to the event and OS
func isModEvent(eventOp fsnotify.Op) bool {
	if runtime.GOOS == "darwin" {
		// depending on the version of Mac OSX, modifying a file will send either a RENAME or a WRITE event
		if eventOp == fsnotify.Rename || eventOp == fsnotify.Write {
			return true
		}
	}

	if runtime.GOOS == "linux" {
		// on linux, modifying a file will send CHMOD events
		if eventOp == fsnotify.Chmod {
			return true
		}
	}

	if runtime.GOOS == "windows" {
		if eventOp == fsnotify.Write {
			return true
		}
	}

	return false
}
