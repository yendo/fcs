package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcs/test"
)

const testNotesFile = "test/test_fcnotes.md"

func openTestNotesFile(t *testing.T) *os.File {
	t.Helper()

	fd, err := os.Open(testNotesFile)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = fd.Close()
		require.NoError(t, err)
	})

	return fd
}

func setCommandLineFlag(t *testing.T, f string) {
	t.Helper()

	err := flag.CommandLine.Set(f, "true")
	require.NoError(t, err)

	t.Cleanup(func() {
		err = flag.CommandLine.Set(f, "false")
		require.NoError(t, err)
	})
}

func TestPrintTitles(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	fd := openTestNotesFile(t)
	printTitles(&buf, fd)

	assert.Equal(t, test.GetExpectedTitles(), buf.String())
}

func TestPrintContents(t *testing.T) {
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

			fd := openTestNotesFile(t)
			printContents(&buf, fd, strings.TrimLeft(tc.title, "# "))
			assert.Equal(t, tc.contents, buf.String())
		})
	}
}

func TestPrintFirstURL(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	fd := openTestNotesFile(t)
	printFirstURL(&buf, fd, "URL")
	assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
}

func TestPrintFirstCmdLine(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	fd := openTestNotesFile(t)
	printFirstCmdLine(&buf, fd, "command-line")
	assert.Equal(t, "ls -l | nl\n", buf.String())
}

func TestIsShellCodeBlockBegin(t *testing.T) {
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
			assert.Equal(t, tc.result, isShellCodeBlockBegin(tc.fence))
		})
	}
}

func TestGetFcsFile(t *testing.T) {
	t.Run("set from environment variable", func(t *testing.T) {
		expectedFileName := "test_file_from_env.md"
		t.Setenv("FCS_NOTES_FILE", expectedFileName)

		fileName, err := getNotesFile()
		assert.NoError(t, err)
		assert.Equal(t, expectedFileName, fileName)
	})

	t.Run("default filename", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")

		home, err := os.UserHomeDir()
		require.NoError(t, err)

		filename, err := getNotesFile()
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(home, "fcnotes.md"), filename)
	})
}

func TestRun(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", testNotesFile)

	oldArgs := os.Args

	t.Cleanup(func() {
		os.Args = oldArgs
	})

	t.Run("cannot access UserHomeDir", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		t.Setenv("HOME", "")

		var buf bytes.Buffer
		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot access user home directory: $HOME is not defined")
		assert.Empty(t, buf.String())
	})

	t.Run("cannot access notes file", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		t.Setenv("HOME", "no_exits")

		var buf bytes.Buffer
		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot access notes file: open no_exits/fcnotes.md: no such file or directory")
		assert.Empty(t, buf.String())
	})

	t.Run("with version flag", func(t *testing.T) {
		setCommandLineFlag(t, "v")

		os.Args = []string{"fcs-cli", "-v"}

		var buf bytes.Buffer
		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, true, *showVersion)
		assert.Equal(t, version+"\n", buf.String())
	})

	t.Run("with url flag", func(t *testing.T) {
		setCommandLineFlag(t, "u")

		t.Run("with no args", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-u"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Equal(t, true, *showURL)
			assert.Equal(t, "", buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-u", "URL"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, true, *showURL)
			assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
		})
	})

	t.Run("with cmd flag", func(t *testing.T) {
		setCommandLineFlag(t, "c")

		t.Run("with no args", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-c"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Equal(t, true, *showCmd)
			assert.Equal(t, "", buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-c", "command-line"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, true, *showCmd)
			assert.Equal(t, "ls -l | nl\n", buf.String())
		})

		t.Run("with a arg has $", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-c", "command-line with $"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, true, *showCmd)
			assert.Equal(t, "date\n", buf.String())
		})
	})

	t.Run("without args", func(t *testing.T) {
		os.Args = []string{"fcs-cli"}

		var buf bytes.Buffer
		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, test.GetExpectedTitles(), buf.String())
	})

	t.Run("with an arg", func(t *testing.T) {
		os.Args = []string{"fcs-cli", "title"}

		var buf bytes.Buffer
		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "# title\n\ncontents\n", buf.String())
	})

	t.Run("with two args", func(t *testing.T) {
		os.Args = []string{"fcs-cli", "title", "other"}

		var buf bytes.Buffer
		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Equal(t, "", buf.String())
	})
}
