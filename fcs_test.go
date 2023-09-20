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

	file := test.OpenTestNotesFile(t)
	fcs.WriteTitles(&buf, file)

	assert.Equal(t, test.GetExpectedTitles(), buf.String())
}

func TestWriteContents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title    string
		contents string
	}{
		{"# title", "# title\n\n" + "contents\n"},
		{"# long title one", "# long title one\n\n" + "line one\nline two\n"},
		{"# title has regular expression meta chars $", "# title has regular expression meta chars $\n\n" + "line\n"},
		{"# contents have blank lines", "# contents have blank lines\n\n" + "1st line\n\n2nd line\n"},
		{"# same title", "# same title\n\n1st\n\n" + "# same title\n\n2nd\n\n" + "# same title\n\n3rd\n"},
		{"## other heading level", "## other heading level\n\n" + "contents\n"},
		{"# title has trailing spaces  ", "# title has trailing spaces  \n\n" + "The contents have trailing spaces.  \n"},
		{"#", ""},
		{"# no contents", "# no contents\n"},
		{"#no_space_title", ""},
		{"#   spaces before title", "#   spaces before title\n\n" + "line\n"},
		{"# fenced code block", "# fenced code block\n\n" + "```\n" + "# fenced heading\n" + "```\n"},
		{"# URL", "# URL\n\n" + "fcs: http://github.com/yendo/fcs/\n" + "github: http://github.com/\n"},
		{"# command-line", "# command-line\n\n" + "```sh\n" + "ls -l | nl\n" + "```\n"},
		{"# no blank line between title and contents", "# no blank line between title and contents\n" + "contents\n"},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			file := test.OpenTestNotesFile(t)
			fcs.WriteContents(&buf, file, strings.TrimLeft(tc.title, "# "))

			assert.Equal(t, tc.contents, buf.String())
		})
	}
}

func TestWriteFirstURL(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	file := test.OpenTestNotesFile(t)
	fcs.WriteFirstURL(&buf, file, "URL")

	assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
}

func TestWriteFirstCmdLine(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	file := test.OpenTestNotesFile(t)
	fcs.WriteFirstCmdLine(&buf, file, "command-line")

	assert.Equal(t, "ls -l | nl\n", buf.String())
}

func TestIsShellCodeBlockBegin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fence  string
		result bool
	}{
		{fence: "```shell", result: true},
		{fence: "``` shell", result: true},
		{fence: "````shell", result: false},
		{fence: "```sh", result: true},
		{fence: "```shell-script", result: true},
		{fence: "```bash", result: true},
		{fence: "```zsh", result: true},
		{fence: "```powershell", result: true},
		{fence: "```posh", result: true},
		{fence: "```pwsh", result: true},
		{fence: "```shellsession", result: true},
		{fence: "```bash session", result: true},
		{fence: "```console", result: true},
		{fence: "```", result: false},
		{fence: "```go", result: false},
		{fence: "```sharp", result: false},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.fence, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.result, fcs.IsShellCodeBlockBegin(tc.fence))
		})
	}
}

func TestWriteNoteLocation(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	file := test.OpenTestNotesFile(t)
	fcs.WriteNoteLocation(&buf, file, "URL")

	assert.Equal(t, fmt.Sprintf("\"%s\" 65\n", test.TestNotesFile), buf.String())
}

func TestGetFcsFile(t *testing.T) {
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
		assert.Equal(t, filepath.Join(home, "fcnotes.md"), filename)
	})
}
