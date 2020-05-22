package main

import (
	"flag"
	"fmt"
)

// flagExt is a struct type that holds a pointer to a boolean value of a flag, and the extension name corresponding to that flag
type flagExt struct {
	extension string
	flagVal   *bool
}

func fileExtensionsToBeWatched() []string {
	// names is a string slice of the file exten
	extensionNames := []string{"js", "html", "css", "hbs", "pug"}

	// flagExtValues is a slice of flagExt values
	flagExtValues := make([]flagExt, 0, len(extensionNames))

	// add a -go flag with a default of true, it is special since we always want to watch go files
	flagExtValues = append(flagExtValues, flagExt{".go", flag.Bool("go", true, "whether or not .go files should be watched for changes")})

	for _, name := range extensionNames {
		// create new flags for every extensionName specified in the above slice
		flagExtValues = append(flagExtValues, flagExt{fmt.Sprintf(".%s", name), flag.Bool(name, false, fmt.Sprintf("whether or not .%s files should be watched for changes", name))})
	}

	flag.Parse()

	var fileExtensionsToWatch []string
	for _, val := range flagExtValues {
		// if a flag has been set to true, add the file extension to fileExtensionsToWatch
		if *val.flagVal {
			fileExtensionsToWatch = append(fileExtensionsToWatch, val.extension)
		}
	}

	return fileExtensionsToWatch
}
