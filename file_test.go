package main

import (
	"testing"
)

func TestFileShouldBeWatched(t *testing.T) {
	if fileShouldBeWatched("", []string{".go", ".html", ".css"}) {
		t.Error(`empty string file name should not be watched`)
	}

	if !fileShouldBeWatched("cmd/dir/foo/bar/baz/main.go", []string{".go", ".html", ".css"}) {
		t.Error("nested file cmd/dir/foo/bar/baz/main.go not watched")
	}

}

func TestFilterOutTest(t *testing.T) {
	if len(filterOutTests([]string{})) != 0 {
		t.Errorf("filterOutTests([]string{}) returned a new array with more elements: %#v", filterOutTests([]string{}))
	}

	if !stringSliceEquals(filterOutTests([]string{"main.go", "other.go", "main_test.go", "other_test.go"}), []string{"main.go", "other.go"}) {
		// fmt.Println(filterOutTests([]string{"main.go", "other.go", "main_test.go", "other_test.go"})[0])
		t.Errorf("tests not removed\n %#v\n", filterOutTests([]string{"main.go", "other.go", "main_test.go", "other_test.go"}))

	}
}

// stringSliceEquals compares two string slices, s and o
// returns true if the slices are the same length and have all of the same values,
// returns false otherwise
func stringSliceEquals(s, o []string) bool {
	if len(s) != len(o) {
		return false
	}

	for i, v := range s {
		if v != o[i] {
			return false
		}
	}
	return true
}
