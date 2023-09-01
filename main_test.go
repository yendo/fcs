package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintTitles(t *testing.T) {
	var buf bytes.Buffer

	fileName := "test/test_memo.md"
	print_titles(&buf, fileName)

	assert.Equal(t, `memo
title1
long title one
contents have blank lines
same title
other title
same title

no contents
no_space_title
last
`, buf.String())
}

func TestPrintContents(t *testing.T) {
	fileName := "test/test_memo.md"

	tests := []struct {
		title    string
		contents string
	}{
		{"## title1", "## title1\n\n" + "contents1\n"},
		{"## long title one", "## long title one\n\n" + "one line1\none line2\n"},
		{"## contents have blank lines", "## contents have blank lines\n\n" + "1st line\n\n2nd line\n"},
		{"## same title", "## same title\n\n" + "contents 1\n\n" + "## same title\n\n" + "contents 2\n"},
		{"## other title", "## other title\n\n" + "other contents\n"},
		{"##", ""},
		{"## no contents", "## no contents\n"},
		{"##no_space_title", ""},
		{"## last", "## last\n\n" + "last\n"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer

		print_contents(&buf, fileName, strings.TrimLeft(tt.title, "# "))
		assert.Equal(t, tt.contents, buf.String())
	}
}
