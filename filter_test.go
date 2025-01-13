package fcqs_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcqs"
)

func TestFilterWriter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input        string
		output       string
		isRemoveHead bool
	}{
		{"a\nb\n", "a\nb\n", false},
		{"\na\nb\n\n", "a\nb\n", false},
		{"a\n\n\nb\n", "a\n\nb\n", false},
		{"a\n\n\n\nb\n", "a\n\nb\n", false},
		{"# a\n\n\n\nb\n", "# a\n\nb\n", false},

		{"a\nb\n", "b\n", true},
		{"\na\nb\n\n", "a\nb\n", true},
		{"a\n\n\nb\n", "b\n", true},
		{"a\n\n\n\nb\n", "b\n", true},
		{"# a\n\n\n\nb\n", "b\n", true},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			f := fcqs.ExportNewFilter(&buf, tc.isRemoveHead)

			file := strings.NewReader(tc.input)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				fcqs.ExportFilterWrite(&f, line)
			}
			err := scanner.Err()
			require.NoError(t, err)

			fcqs.ExportFilterWrite(&f, "")

			assert.Equal(t, tc.output, buf.String())
		})
	}
}
