package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	notesFile         = "testdata/test_fcnotes.md"
	shellBlockFile    = "testdata/test_shellblock.md"
	locationFile      = "testdata/test_location.md"
	locationExtraFile = "testdata/test_location_extra.md"
)

var (
	NotesFile         string
	ShellBlockFile    string
	LocationFile      string
	LocationExtraFile string
)

func init() {
	NotesFile = fullPath(notesFile)
	ShellBlockFile = fullPath(shellBlockFile)
	LocationFile = fullPath(locationFile)
	LocationExtraFile = fullPath(locationExtraFile)
}

// OpenTestNotesFile opens a test notes file.
func OpenTestNotesFile(t *testing.T, filename string) *os.File {
	t.Helper()

	file, err := os.Open(filename)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := file.Close()
		require.NoError(t, err)
	})

	return file
}

// fullPath returns a full path of test note file.
func fullPath(filename string) string {
	_, thisFileName, _, _ := runtime.Caller(0)

	return filepath.Join(path.Dir(thisFileName), filename)
}

// ExpectedTitles returns the titles of the test notes.
func ExpectedTitles() string {
	return `title
Long title and contents have lines
Regular expression meta chars in the title are ignored $
Consecutive blank lines are combined into a single line
same title
Heading levels and structures are ignored
Trailing spaces in the title are ignored
Spaces before the title are ignored
Headings in fenced code blocks are ignored
There can be no blank line
Titles without a space after the # are not recognized
URL
command-line
command-line with $
`
}
