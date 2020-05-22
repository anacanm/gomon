package main

import "testing"

func TestFileShouldBeWatched(t *testing.T) {
	if fileShouldBeWatched("", []string{".go", ".html", ".css"}) {
		t.Error(`empty string file name should not be watched`)
	}

	if !fileShouldBeWatched("cmd/dir/foo/bar/baz/main.go", []string{".go", ".html", ".css"}) {
		t.Error("nested file cmd/dir/foo/bar/baz/main.go not watched")
	}

}
