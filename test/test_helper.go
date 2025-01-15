package test

import (
	_ "embed"
	"path"
	"path/filepath"
	"runtime"
)

const (
	notesFile         = "testdata/test_fcnotes.md"
	shellBlockFile    = "testdata/test_shellblock.md"
	locationFile      = "testdata/test_location.md"
	locationExtraFile = "testdata/test_location_extra.md"
)

var (
	NotesFile         = fullPath(notesFile)
	ShellBlockFile    = fullPath(shellBlockFile)
	LocationFile      = fullPath(locationFile)
	LocationExtraFile = fullPath(locationExtraFile)
)

// fullPath returns a full path of test note file.
func fullPath(filename string) string {
	_, thisFileName, _, _ := runtime.Caller(0)

	return filepath.Join(path.Dir(thisFileName), filename)
}

// ExpectedTitles has titles of the test notes.
//
//go:embed testdata/expected_titles.txt
var ExpectedTitles string
