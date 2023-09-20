package test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const TestNotesFile = "test/test_fcnotes.md"

func OpenTestNotesFile(t *testing.T) *os.File {
	t.Helper()

	fd, err := os.Open(TestNotesFile)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := fd.Close()
		require.NoError(t, err)
	})

	return fd
}

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
fenced code block
URL
command-line
command-line with $
no blank line between title and contents
`

	return strings.Replace(titles, "trailing spaces", "trailing spaces  ", 1)
}
