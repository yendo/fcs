package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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
	titles := `title
long title one
title has regular expression meta chars $
contents have blank lines
same title
other heading level
title has trailing spaces
no contents
no contents2
no_space_title
spaces before title
fenced code block
URL
command-line
command-line with $
no blank line between title and contents
`

	return strings.Replace(titles, "trailing spaces", "trailing spaces  ", 1)
}
