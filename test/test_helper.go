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
	TestNotesFile      = "testdata/test_fcnotes.md"
	TestShellBlockFile = "testdata/test_shellblock.md"
	TestLocationFile   = "testdata/test_location.md"
)

// OpenTestNotesFile opens a test notes file.
func OpenTestNotesFile(t *testing.T, filename string) *os.File {
	t.Helper()

	file, err := os.Open(GetTestDataFullPath(filename))
	require.NoError(t, err)

	t.Cleanup(func() {
		err := file.Close()
		require.NoError(t, err)
	})

	return file
}

func GetTestDataFullPath(filename string) string {
	_, thisFileName, _, _ := runtime.Caller(0)

	return filepath.Join(path.Dir(thisFileName), filename)
}

// GetExpectedTitles returns the titles of the test notes.
func GetExpectedTitles() string {
	return `title
Long title and contents have lines
Regular expression meta chars in the title are ignored $
Consecutive blank lines are combined into a single line
same title
Heading levels and structures are ignored
Trailing spaces in the title are ignored
Notes without content output the title only
Spaces before the title are ignored
Headings in fenced code blocks are ignored
There can be no blank line
Titles without a space after the # are not recognized
URL
command-line
command-line with $
`
}
