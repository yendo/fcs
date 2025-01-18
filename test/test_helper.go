package test

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// MultiFiles returns file names concatenated with PathListSeparator for FCQS_NOTES_FILE.
func MultiFiles(files ...string) string {
	return strings.Join(files, string(os.PathListSeparator))
}

// fullPath returns a full path of test note file.
func fullPath(filename string) string {
	_, thisFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("fail to get a test file path")
	}

	return filepath.Join(filepath.Dir(thisFileName), filename)
}

// ExpectedTitles has titles of the test notes.
//
//go:embed testdata/expected_titles.txt
var ExpectedTitles string
