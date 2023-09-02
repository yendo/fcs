package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintTitles(t *testing.T) {
	var buf bytes.Buffer

	fileName := "test/test_fcnotes.md"
	fd, err := os.Open(fileName)
	require.NoError(t, err)
	defer fd.Close()

	printTitles(&buf, fd)

	assert.Equal(t, `notes
title
long title one
contents have blank lines
same title
other title
no contents
no contents2
no_space_title
no blank line between title and contents
`, buf.String())
}

func TestPrintContents(t *testing.T) {
	fileName := "test/test_fcnotes.md"

	tests := []struct {
		title    string
		contents string
	}{
		{"## title", "## title\n\n" + "contents\n"},
		{"## long title one", "## long title one\n\n" + "line one\nline two\n"},
		{"## contents have blank lines", "## contents have blank lines\n\n" + "1st line\n\n2nd line\n"},
		{"## same title", "## same title\n\ncontents 1\n\n" + "## same title\n\ncontents 2\n\n" + "## same title\n\ncontents 3\n"},
		{"## other title", "## other title\n\n" + "other contents\n"},
		{"##", ""},
		{"## no contents", "## no contents\n"},
		{"##no_space_title", ""},
		{"## no blank line between title and contents", "## no blank line between title and contents\n" + "contents\n"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		fd, err := os.Open(fileName)
		require.NoError(t, err)
		defer fd.Close()

		printContents(&buf, fd, strings.TrimLeft(tt.title, "# "))
		assert.Equal(t, tt.contents, buf.String())
	}
}

func TestRun(t *testing.T) {

	t.Run("no args", func(t *testing.T) {
		var buf bytes.Buffer
		run(&buf)

		assert.Equal(t, `notes
title
long title one
contents have blank lines
same title
other title
no contents
no contents2
no_space_title
no blank line between title and contents
`, buf.String())
	})

	t.Run("with a arg", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"fcs-cli", "title"}

		var buf bytes.Buffer
		run(&buf)

		assert.Equal(t, "## title\n\ncontents\n", buf.String())
	})

}
