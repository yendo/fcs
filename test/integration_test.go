package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const defaultNotesFile = "fcnotes.md"

type testCmd struct {
	cmd    *exec.Cmd
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (c *testCmd) run() error {
	c.cmd.Stdout = &c.stdout
	c.cmd.Stderr = &c.stderr

	return c.cmd.Run()
}

func newTestCmd(args ...string) *testCmd {
	cmd := exec.Command("./fcqs-cli", args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=../coverdir")

	return &testCmd{cmd: cmd}
}

func TestCmdSuccess(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", NotesFile)
	t.Setenv("FCQS_NOTES_FILES", "")

	tests := []struct {
		title   string
		options []string
		stdout  string
	}{
		{
			title:   "with version flag",
			options: []string{"-v"},
			stdout:  "0.0.0-test\n",
		},
		{
			title:   "with url flag and an arg",
			options: []string{"-u", "URL"},
			stdout:  "http://github.com/yendo/fcqs/\n",
		},
		{
			title:   "with url flag and an empty arg",
			options: []string{"-u", ""},
			stdout:  "",
		},
		{
			title:   "with cmd flag and an arg",
			options: []string{"-c", "more command-line blocks"},
			stdout:  "ls -l | nl\n",
		},
		{
			title:   "with cmd flag and an empty arg",
			options: []string{"-c", ""},
			stdout:  "",
		},
		{
			title:   "with cmd flag and an arg not for code block",
			options: []string{"-c", "Headings in fenced code blocks are ignored"},
			stdout:  "",
		},
		{
			title:   "with location flag and an arg",
			options: []string{"-l", "title"},
			stdout:  fmt.Sprintf("%q 1\n", NotesFile),
		},
		{
			title:   "with location flag and an empty arg",
			options: []string{"-l", ""},
			stdout:  "",
		},
		{
			title:   "without args",
			options: []string{},
			stdout:  ExpectedTitles,
		},
		{
			title:   "with an empty arg",
			options: []string{""},
			stdout:  "",
		},
		{
			title:   "with an arg of spaces",
			options: []string{"  "},
			stdout:  "",
		},
		{
			title:   "with an arg",
			options: []string{"There can be no blank line"},
			stdout:  "# There can be no blank line\ncontents\n", // next line is only "#".
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			cmd := newTestCmd(tc.options...)
			err := cmd.run()

			require.NoError(t, err)
			assert.Equal(t, tc.stdout, cmd.stdout.String())
			assert.Empty(t, cmd.stderr.String())
		})
	}
}

func TestCmdWriteContentsWithoutTitle(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", NotesFile)
	t.Setenv("FCQS_NOTES_FILES", "")

	cmd := newTestCmd("-t", "There can be no blank line")
	err := cmd.run()

	require.NoError(t, err)
	assert.Equal(t, "contents\n", cmd.stdout.String())
	assert.Empty(t, cmd.stderr.String())
}

func TestCmdFail(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", NotesFile)
	t.Setenv("FCQS_NOTES_FILES", "")

	tests := []struct {
		title   string
		options []string
		stderr  string
	}{
		{
			title:   "with url flag and no arg",
			options: []string{"-u"},
			stderr:  "invalid number of arguments\n",
		},
		{
			title:   "with cmd flag and no arg",
			options: []string{"-c"},
			stderr:  "invalid number of arguments\n",
		},
		{
			title:   "with location flag and no arg",
			options: []string{"-l"},
			stderr:  "invalid number of arguments\n",
		},
		{
			title:   "with two args",
			options: []string{"title", "other"},
			stderr:  "invalid number of arguments\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			cmd := newTestCmd(tc.options...)
			err := cmd.run()

			require.Error(t, err)
			assert.Empty(t, cmd.stdout.String())
			assert.Equal(t, tc.stderr, cmd.stderr.String())
		})
	}
}

func TestCmdNotesLocation(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", LocationFile)
	t.Setenv("FCQS_NOTES_FILES", "")

	cmd := newTestCmd("-l", "5th Line")
	err := cmd.run()

	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%q 5\n", LocationFile), cmd.stdout.String())
	assert.Empty(t, cmd.stderr.String())
}

func TestCmdMultiFiles(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", MultiFiles(LocationFile, LocationExtraFile))
	t.Setenv("FCQS_NOTES_FILES", "")

	t.Run("show titles", func(t *testing.T) {
		cmd := newTestCmd()
		err := cmd.run()

		require.NoError(t, err)
		assert.Equal(t, "location test data\n5th Line\nother 5th Line\n9th Line\n", cmd.stdout.String())
		assert.Empty(t, cmd.stderr.String())
	})

	t.Run("show contents", func(t *testing.T) {
		cmd := newTestCmd("9th Line")
		err := cmd.run()

		require.NoError(t, err)
		assert.Equal(t, "# 9th Line\n\nDo not chang the 9th line.\n", cmd.stdout.String())
		assert.Empty(t, cmd.stderr.String())
	})

	t.Run("show location", func(t *testing.T) {
		cmd := newTestCmd("-l", "9th Line")
		err := cmd.run()

		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%q 9\n", LocationExtraFile), cmd.stdout.String())
		assert.Empty(t, cmd.stderr.String())
	})

	t.Run("file error", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", MultiFiles(LocationFile, "invalid_file"))
		t.Setenv("FCQS_NOTES_FILES", "")

		cmd := newTestCmd()
		err := cmd.run()

		require.Error(t, err)
		assert.Equal(t, "notes file: open invalid_file: no such file or directory\n", cmd.stderr.String())
		assert.Empty(t, cmd.stdout.String())
	})
}

func TestUserHomeDirNotExists(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", "")
	t.Setenv("FCQS_NOTES_FILES", "")
	t.Setenv("HOME", "")

	cmd := newTestCmd()
	err := cmd.run()

	require.Error(t, err)
	assert.Empty(t, cmd.stdout.String())
	assert.Equal(t, "notes file name: user home directory: $HOME is not defined\n", cmd.stderr.String())
}

func TestNotesNotExists(t *testing.T) {
	t.Setenv("FCQS_NOTES_FILE", "not_exists")
	t.Setenv("FCQS_NOTES_FILES", "")

	cmd := newTestCmd()
	err := cmd.run()

	require.Error(t, err)
	assert.Empty(t, cmd.stdout.String())
	assert.Equal(t, "notes file: open not_exists: no such file or directory\n", cmd.stderr.String())
}

func TestDefaultNoteExists(t *testing.T) {
	// Arrange
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("FCQS_NOTES_FILE", "")
	t.Setenv("FCQS_NOTES_FILES", "")

	file := filepath.Join(tempHome, defaultNotesFile)
	err := os.WriteFile(file, []byte("# title\ncontents\n"), 0o600)
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	require.Equal(t, file, filepath.Join(home, defaultNotesFile))

	// Act
	cmd := newTestCmd()
	err = cmd.run()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "title\n", cmd.stdout.String())
	assert.Empty(t, cmd.stderr.String())
}

func TestBashScript(t *testing.T) {
	t.Parallel()

	// Arrange
	fileName := "../shell.bash"
	data, err := os.ReadFile(fileName)
	require.NoError(t, err)
	expected := string(data)

	// Act
	cmd := newTestCmd("--bash")
	err = cmd.run()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, cmd.stdout.String())
	assert.Empty(t, cmd.stderr.String())
}
