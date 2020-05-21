package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	// begin walking at the root of the directory specified, adding a fsnotify watcher to each file
	// TODO: make root path specifiable by user input
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
				fmt.Printf("%#v\n", event)

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
