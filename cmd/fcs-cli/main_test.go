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

func TestRun(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", test.GetTestDataFullPath(test.TestNotesFile))

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

	t.Run("with version flag", func(t *testing.T) {
		setCommandLineFlag(t, "v")
		var buf bytes.Buffer
		os.Args = []string{"fcs-cli", "-v"}

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, version+"\n", buf.String())
	})

	t.Run("with url flag", func(t *testing.T) {
		setCommandLineFlag(t, "u")

		t.Run("with no args", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-u"}

			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Empty(t, buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-u", "URL"}

			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, "http://github.com/yendo/fcs/\n", buf.String())
		})
	})

	t.Run("with cmd flag", func(t *testing.T) {
		setCommandLineFlag(t, "c")

		t.Run("with no args", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-c"}

			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Empty(t, buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-c", "command-line"}

			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, "ls -l | nl\n", buf.String())
		})

		t.Run("with a arg has $", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-c", "command-line with $"}

			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, "date\n", buf.String())
		})
	})

	t.Run("with location flag", func(t *testing.T) {
		testFileName := test.GetTestDataFullPath(test.TestLocationFile)
		t.Setenv("FCS_NOTES_FILE", testFileName)
		setCommandLineFlag(t, "l")

		t.Run("with no args", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-l"}

			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Empty(t, buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			var buf bytes.Buffer
			os.Args = []string{"fcs-cli", "-l", "5th Line"}

			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("%q 5\n", testFileName), buf.String())
		})
	})

	t.Run("without args", func(t *testing.T) {
		var buf bytes.Buffer
		os.Args = []string{"fcs-cli"}

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, test.GetExpectedTitles(), buf.String())
	})

	t.Run("with an arg", func(t *testing.T) {
		var buf bytes.Buffer
		os.Args = []string{"fcs-cli", "title"}

		err := run(&buf)

		assert.NoError(t, err)
		assert.Equal(t, "# title\n\ncontents\n", buf.String())
	})

	t.Run("with two args", func(t *testing.T) {
		var buf bytes.Buffer
		os.Args = []string{"fcs-cli", "title", "other"}

		err := run(&buf)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid number of arguments")
		assert.Empty(t, buf.String())
	})
}
