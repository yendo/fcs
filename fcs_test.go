package fcs_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcs"
	"github.com/yendo/fcs/test"
)

func TestWriteTitles(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	file := test.OpenTestNotesFile(t, test.TestNotesFile)
	fcs.WriteTitles(&buf, file)

	assert.Equal(t, test.GetExpectedTitles(), buf.String())
}

func TestWriteContents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title    string
		contents string
	}{
		{"# title\n", "contents\n"},
		{"# Long title and contents have lines\n", "line 1\n\nline 2\n"},
		{"# Regular expression meta chars in the title are ignored $\n", "contents\n"},
		{"# Consecutive blank lines are combined into a single line\n", "line 1\n\nline 2\n"},
		{"# same title\n", "Contents with the same title are combined into one.\n\n" + "# same title\n\n2nd\n\n" + "# same title\n\n3rd\n"},
		{"## Heading levels and structures are ignored\n", "contents\n"},
		{"# Trailing spaces in the title are ignored  \n", "The contents have trailing spaces.  \n"},
		{"# Notes without content output the title only", ""},
		{"#   Spaces before the title are ignored\n", "contents\n"},
		{"# Headings in fenced code blocks are ignored\n", "```\n" + "# fenced heading\n" + "```\n"},
		{"# There can be no blank line", "contents\n"},
		{"# URL\n", "fcs: http://github.com/yendo/fcs/\n" + "github: http://github.com/\n"},
		{"# command-line\n", "```sh\n" + "ls -l | nl\n" + "```\n"},
		{"# command-line with $\n", "```console\n" + "$ date\n" + "```\n"},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			file := test.OpenTestNotesFile(t, test.TestNotesFile)
			var buf bytes.Buffer
			title := strings.TrimRight(strings.TrimLeft(tc.title, "# "), "\n")

			fcs.WriteContents(&buf, file, title)

			assert.Equal(t, tc.title+"\n"+tc.contents, buf.String())
		})
	}
}

func TestWriteNoContents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		title string
	}{
		{
			desc:  "Titles without strings are not recognized.",
			title: "#",
		},
		{
			desc:  "Titles without a space after the `#` are not recognized",
			title: "#no_space_title",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			file := test.OpenTestNotesFile(t, test.TestNotesFile)

			fcs.WriteContents(&buf, file, strings.TrimLeft(tc.title, "#"))

			assert.Empty(t, buf.String())
		})
	}
}

func TestWriteFirstURL(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	file := test.OpenTestNotesFile(t, test.TestNotesFile)

	fcs.WriteFirstURL(&buf, file, "URL")

	assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
}

func TestWriteFirstCmdLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title  string
		output bool
	}{
		{"shell 1", true},
		{"shell 2", true},
		{"shell 3", true},
		{"sh", true},
		{"shell-script", true},
		{"bash", true},
		{"zsh", true},
		{"powershell", true},
		{"posh", true},
		{"pwsh", true},
		{"shellsession", true},
		{"bash session", true},
		{"console", true},
		{"go", false},
		{"no identifier", false},
		{"other identifier", false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			file := test.OpenTestNotesFile(t, test.TestShellBlockFile)

			fcs.WriteFirstCmdLine(&buf, file, tc.title)

			expected := map[bool]string{true: "ls -l | nl\n", false: ""}
			assert.Equal(t, expected[tc.output], buf.String())
		})
	}
}

func TestWriteNoteLocation(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	testFile := test.OpenTestNotesFile(t, test.TestLocationFile)
	fcs.WriteNoteLocation(&buf, testFile, "5th Line")

	assert.Equal(t, fmt.Sprintf("%q 5\n", testFile.Name()), buf.String())
}

func TestGetFcsFile(t *testing.T) {
	t.Run("cannot access user home directory", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		t.Setenv("HOME", "")

		fileName, err := fcs.GetNotesFileName()

		assert.Empty(t, fileName)
		assert.Error(t, err)
		assert.EqualError(t, err, "cannot access user home directory: $HOME is not defined")
	})

	t.Run("set from environment variable", func(t *testing.T) {
		expectedFileName := "test_file_from_env.md"
		t.Setenv("FCS_NOTES_FILE", expectedFileName)

		fileName, err := fcs.GetNotesFileName()

		assert.NoError(t, err)
		assert.Equal(t, expectedFileName, fileName)
	})

	t.Run("default filename", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		home, err := os.UserHomeDir()
		require.NoError(t, err)

		filename, err := fcs.GetNotesFileName()

		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(home, fcs.DefaultNotesFile), filename)
	})
}
