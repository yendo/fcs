package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stdBuf struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (b *stdBuf) newTestCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("./fcs-cli", args...)
	cmd.Env = append(os.Environ(), "GOCOVERDIR=../coverdir")
	cmd.Stdout = &b.stdout
	cmd.Stderr = &b.stderr

	return cmd
}

func TestCmd(t *testing.T) {
	tests := []struct {
		title   string
		options []string
		err     bool
		stdout  string
		stderr  string
	}{
		{
			title:   "with version flag",
			options: []string{"-v"},
			err:     false,
			stdout:  "0.0.0-test\n",
			stderr:  "",
		},
		{
			title:   "with url flag and no arg",
			options: []string{"-u"},
			err:     true,
			stdout:  "",
			stderr:  "invalid number of arguments\n",
		},
		{
			title:   "with url flag and an arg",
			options: []string{"-u", "URL"},
			err:     false,
			stdout:  "http://github.com/yendo/fcs/\n",
			stderr:  "",
		},
		{
			title:   "without args",
			options: []string{},
			err:     false,
			stdout:  GetExpectedTitles(),
			stderr:  "",
		},
		{
			title:   "with an arg",
			options: []string{"title"},
			err:     false,
			stdout:  "# title\n\ncontents\n",
			stderr:  "",
		},
		{
			title:   "with two args",
			options: []string{"title", "other"},
			err:     true,
			stdout:  "",
			stderr:  "invalid number of arguments\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			t.Setenv("FCS_NOTES_FILE", "test_fcnotes.md")

			buf := &stdBuf{}
			cmd := buf.newTestCmd(tt.options...)
			err := cmd.Run()

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.stdout, buf.stdout.String())
			assert.Equal(t, tt.stderr, buf.stderr.String())
		})
	}
}

func TestNotesNotExists(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", "not_exists")

	buf := &stdBuf{}
	cmd := buf.newTestCmd()
	err := cmd.Run()

	assert.Error(t, err)
	assert.Empty(t, buf.stdout.String())
	assert.Equal(t, "open not_exists: no such file or directory\n", buf.stderr.String())
}

func TestDefaultNoteExists(t *testing.T) {
	t.Setenv("FCS_NOTES_FILE", "")

	home, err := os.UserHomeDir()
	require.NoError(t, err)
	if _, err := os.Stat(filepath.Join(home, "fcnotes.md")); err != nil {
		t.Skip("the default fcnotes.md does not exist")
	}

	buf := &stdBuf{}
	cmd := buf.newTestCmd()
	err = cmd.Run()

	assert.NoError(t, err)
	assert.Empty(t, buf.stderr.String())
}
