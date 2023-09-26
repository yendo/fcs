package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcs"
	"github.com/yendo/fcs/test"
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
	t.Setenv("FCS_NOTES_FILE", test.GetTestDataFullPath(test.TestNotesFile))

	t.Run("without args", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, test.GetExpectedTitles(), buf.String())
	})

	t.Run("with an arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "title"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "# title\n\ncontents\n", buf.String())
	})
}

func TestRunFail(t *testing.T) {
	t.Run("cannot access UserHomeDir", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		t.Setenv("HOME", "")
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot get notes file name: cannot access user home directory: $HOME is not defined")
		assert.Empty(t, buf.String())
	})

	t.Run("cannot access notes file", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", "")
		t.Setenv("HOME", "no_exits")
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("cannot access notes file: open no_exits/%s: no such file or directory", fcs.DefaultNotesFile))
		assert.Empty(t, buf.String())
	})

	t.Run("with two args", func(t *testing.T) {
		t.Setenv("FCS_NOTES_FILE", test.GetTestDataFullPath(test.TestNotesFile))
		setOSArgs(t, []string{"fcs-cli", "title", "other"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})
}

func TestRunWithVersionFlag(t *testing.T) {
	setCommandLineFlag(t, "v")
	setOSArgs(t, []string{"fcs-cli", "-v"})
	var buf bytes.Buffer

	err := run(&buf)

	assert.NoError(t, err)
	assert.Equal(t, version+"\n", buf.String())
}

func TestRunWithURLFlag(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", test.GetTestDataFullPath(test.TestNotesFile))
	setCommandLineFlag(t, "u")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-u"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-u", "URL"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
	})
}

func TestRunWithCmdFlag(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", test.GetTestDataFullPath(test.TestNotesFile))
	setCommandLineFlag(t, "c")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-c"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-c", "command-line"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "ls -l | nl\n", buf.String())
	})

	t.Run("with a arg has $", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-c", "command-line with $"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "date\n", buf.String())
	})
}

func TestRunWithLocationFlag(t *testing.T) {
	testFileName := test.GetTestDataFullPath(test.TestLocationFile)
	t.Setenv("FCS_NOTES_FILE", testFileName)
	setCommandLineFlag(t, "l")

	t.Run("with no args", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-l"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		setOSArgs(t, []string{"fcs-cli", "-l", "5th Line"})
		var buf bytes.Buffer

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%q 5\n", testFileName), buf.String())
	})
}
