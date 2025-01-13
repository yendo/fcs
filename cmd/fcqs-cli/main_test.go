package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcqs"
	"github.com/yendo/fcqs/test"
)

func setCommandLineFlag(t *testing.T, f string) {
	t.Helper()

	err := flag.CommandLine.Set(f, "true")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := flag.CommandLine.Set(f, "false")
		require.NoError(t, err)
	})
}

func setOSArgs(t *testing.T, args []string) {
	t.Helper()

	oldArgs := os.Args
	os.Args = args

	t.Cleanup(func() {
		os.Args = oldArgs
	})
}

func TestRunSuccess(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.NotesFile))

	t.Run("without args", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, test.ExpectedTitles(), buf.String())
	})

	t.Run("with an arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "title"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "# title\n\ncontents\n", buf.String())
	})

	t.Run("with an arg and no_title option", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-t", "title"})
		setCommandLineFlag(t, "notitle")

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "contents\n", buf.String())
	})

	t.Run("with an empty arg and some option", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-u", ""})
		setCommandLineFlag(t, "url")

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Empty(t, buf.String())
	})
}

func TestRunFail(t *testing.T) {
	t.Run("failed to access UserHomeDir", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", "")
		t.Setenv("HOME", "")

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "notes file name: user home directory: $HOME is not defined")
		assert.Empty(t, buf.String())
	})

	t.Run("failed to access notes file", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", "")
		t.Setenv("HOME", "no_exits")

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("notes file: open no_exits/%s: no such file or directory", fcqs.DefaultNotesFile))
		assert.Empty(t, buf.String())
	})

	t.Run("with two args", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.NotesFile))
		setOSArgs(t, []string{"fcqs-cli", "title", "other"})

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})
}

func TestRunMultiFiles(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.LocationFile)+string(os.PathListSeparator)+test.FullPath(test.LocationExtraFile))

	t.Run("show titles", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "location test data\n5th Line\nother 5th Line\n9th Line\n", buf.String())
	})

	t.Run("show contents", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "9th Line"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "# 9th Line\n\nDo not chang the 9th line.\n", buf.String())
	})

	t.Run("show location", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-l", "9th Line"})
		setCommandLineFlag(t, "location")

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%q 9\n", test.FullPath(test.LocationExtraFile)), buf.String())
	})

	t.Run("file error", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.LocationFile)+string(os.PathListSeparator)+"invalid_file")
		setOSArgs(t, []string{"fcqs-cli"})

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "notes file: open invalid_file: no such file or directory")
	})
}

func TestRunWithURLFlag(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.NotesFile))
	setCommandLineFlag(t, "url")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-u"})

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-u", "URL"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "http://github.com/yendo/fcqs/\n", buf.String())
	})
}

func TestRunWithCmdFlag(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", test.FullPath(test.NotesFile))
	setCommandLineFlag(t, "command")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-c"})

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-c", "command-line"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "ls -l | nl\n", buf.String())
	})

	t.Run("with a arg has $", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-c", "command-line with $"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, "date\n", buf.String())
	})
}

func TestRunWithLocationFlag(t *testing.T) {
	testFileName := test.FullPath(test.LocationFile)
	t.Setenv("FCQS_NOTES_FILE", testFileName)
	setCommandLineFlag(t, "location")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-l"})

		var buf bytes.Buffer
		err := run(&buf)

		require.Error(t, err)
		require.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcqs-cli", "-l", "5th Line"})

		var buf bytes.Buffer
		err := run(&buf)

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%q 5\n", testFileName), buf.String())
	})
}

func TestRunWithVersionFlag(t *testing.T) {
	t.Parallel()

	setCommandLineFlag(t, "version")
	setOSArgs(t, []string{"fcqs-cli", "-v"})

	var buf bytes.Buffer
	err := run(&buf)

	require.NoError(t, err)
	assert.Equal(t, version+"\n", buf.String())
}

func TestRunWithBashScriptFlag(t *testing.T) {
	setCommandLineFlag(t, "bash")
	setOSArgs(t, []string{"fcqs-cli", "--bash"})

	fileName := "../../shell.bash"
	expectedData, err := os.ReadFile(fileName)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = run(&buf)

	require.NoError(t, err)
	assert.Equal(t, string(expectedData), buf.String())
}
