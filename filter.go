package fcqs

import (
	"fmt"
	"io"
)

// filter represents filter that eliminates unnecessary lines.
type filter struct {
	w            io.Writer
	isRemoveHead bool
	isDoneTitle  bool
	prevLine     []byte
}

// Write writes filtered lines.
func (f *filter) Write(p []byte) (n int, err error) {
	if f.prevLine == nil {
		// Remove first non blank lines.
		if f.isRemoveHead && !f.isDoneTitle {
			f.isDoneTitle = true
			return
		}

		// Remove first blank lines.
		if len(p) != 0 {
			f.prevLine = make([]byte, len(p), len(p)+1) // cap for LF
			copy(f.prevLine, p)
		}
		return
	}

	// Remove consecutive blank lines
	if len(f.prevLine) == 0 && len(p) == 0 {
		return
	}

	// Add line feed
	f.prevLine = append(f.prevLine, '\n')

	n, err = f.w.Write(f.prevLine)
	if err != nil {
		return
	}
	if n != len(f.prevLine) {
		err = io.ErrShortWrite
		return
	}

	f.prevLine = make([]byte, len(p), len(p)+1) // cap for LF
	copy(f.prevLine, p)

	return len(f.prevLine), nil
}

// Close closes the filter.
func (f *filter) Close() error {
	fmt.Fprint(f, "")
	return nil
}

// newFilter returns a filter.
func newFilter(w io.Writer, isRemoveHead bool) io.WriteCloser {
	return &filter{w: w, isRemoveHead: isRemoveHead}
}
