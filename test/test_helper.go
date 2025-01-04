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
	NotesFile         = "testdata/test_fcnotes.md"
	ShellBlockFile    = "testdata/test_shellblock.md"
	LocationFile      = "testdata/test_location.md"
	LocationExtraFile = "testdata/test_location_extra.md"
)

// OpenTestNotesFile opens a test notes file.
func OpenTestNotesFile(t *testing.T, filename string) *os.File {
	t.Helper()

	file, err := os.Open(FullPath(filename))
	require.NoError(t, err)

	t.Cleanup(func() {
		err := file.Close()
		require.NoError(t, err)
	})

	return file
}

// FullPath returns a full path of test note file.
func FullPath(filename string) string {
	_, thisFileName, _, _ := runtime.Caller(0)

	return filepath.Join(path.Dir(thisFileName), filename)
}

func FileSeparator() string {
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}

	return sep
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
