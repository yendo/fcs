package fcqs_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
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
				_, err := f.Write([]byte(line))
				require.NoError(t, err)
			}
			f.Close()

			err := scanner.Err()
			require.NoError(t, err)
			assert.Equal(t, tc.output, buf.String())
		})
	}
}

// writerMock represents a mock for io.Writer.
type writerMock struct {
	num int
	err error
}

func (w writerMock) Write(_ []byte) (int, error) {
	return w.num, w.err
}

func TestFilterWriterError(t *testing.T) {
	t.Parallel()

	outputText := "second line"

	t.Run("Writer returns write err", func(t *testing.T) {
		t.Parallel()

		writeErr := "write error"
		textLen := len(outputText)
		w := &writerMock{num: textLen, err: errors.New(writeErr)}

		f := fcqs.ExportNewFilter(w, false)

		n, err := f.Write([]byte("first line"))
		require.NoError(t, err)
		assert.Zero(t, n)

		n, err = f.Write([]byte(outputText))
		require.Error(t, err)
		assert.Equal(t, textLen, n)
		assert.EqualError(t, err, writeErr)
	})

	t.Run("Writer returns short write error", func(t *testing.T) {
		t.Parallel()

		textLen := len(outputText) - 1
		w := &writerMock{num: textLen, err: nil}

		f := fcqs.ExportNewFilter(w, false)

		n, err := f.Write([]byte("first line"))
		require.NoError(t, err)
		assert.Zero(t, n)

		n, err = f.Write([]byte(outputText))
		require.Error(t, err)
		require.ErrorIs(t, err, io.ErrShortWrite)
		assert.Equal(t, textLen, n)
	})
}
