package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("with location flag", func(t *testing.T) {
		setCommandLineFlag(t, "l")

		t.Run("with no args", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-l"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.Error(t, err)
			assert.EqualError(t, err, "invalid number of arguments")
			assert.Equal(t, true, *showLoc)
			assert.Equal(t, "", buf.String())
		})

		t.Run("with a arg", func(t *testing.T) {
			os.Args = []string{"fcs-cli", "-l", "command-line"}

			var buf bytes.Buffer
			err := run(&buf)

			assert.NoError(t, err)
			assert.Equal(t, true, *showLoc)
			assert.Equal(t, fmt.Sprintf("%q 70\n", test.GetTestDataFullPath(test.TestNotesFile)), buf.String())
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
