package fcqs

import (
	"fmt"
	"io"
)

type filter struct {
	WriteBuf     io.Writer
	IsRemoveHead bool
	isDoneTitle  bool
	prevLine     *string
}

// write writes filtered lines.
func (f *filter) write(text string) {
	if f.prevLine == nil {
		// Remove first non blank lines.
		if f.IsRemoveHead && !f.isDoneTitle {
			f.isDoneTitle = true
			return
		}

		// Remove first blank lines.
		if text != "" {
			f.prevLine = &text
		}
		return
	}

	// Remove consecutive blank lines
	if *f.prevLine == "" && text == "" {
		return
	}

	fmt.Fprintln(f.WriteBuf, *f.prevLine)
	f.prevLine = &text
}

func newFilter(w io.Writer, isRemoveHead bool) filter {
	return filter{WriteBuf: w, IsRemoveHead: isRemoveHead}
}
