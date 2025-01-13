package value

import (
	"fmt"
	"strings"
)

const atxHeadingChar = "#"

// Title represents a title that does not allow empty titles.
type Title struct {
	value string
}

// String returns a title string.
func (t Title) String() string {
	return t.value
}

// Equal reports whether the title equals to other title.
func (t Title) Equals(other *Title) bool {
	return t.String() == other.String()
}

// NewTitleLine returns title.
func NewTitle(t string) (*Title, error) {
	t = strings.Trim(t, " ")
	if t == "" {
		return nil, fmt.Errorf("title is empty")
	}

	return &Title{value: t}, nil
}

// TitleLine represents a title text line that allows empty titles.
type TitleLine struct {
	title *Title
}

// Title returns a title in the title line.
func (tl TitleLine) Title() Title {
	return *tl.title
}

// HasValidTitle reports whether a title in the title line is valid.
func (tl TitleLine) HasValidTitle() bool {
	return tl.title != nil
}

// EqualTitle reports whether a title equals to a title in the title line.
func (tl TitleLine) EqualTitle(title *Title) bool {
	return tl.title != nil && title.Equals(tl.title)
}

// NewTitleLine returns title line.
func NewTitleLine(tl string) (*TitleLine, bool) {
	if !isTitleLine(tl) {
		return nil, false
	}

	titleStr := strings.Trim(tl, atxHeadingChar+" ")
	title, err := NewTitle(titleStr)
	if err != nil {
		return &TitleLine{title: nil}, true
	}

	return &TitleLine{title: title}, true
}

// isTitleLine returns if the line is title line.
func isTitleLine(line string) bool {
	// Title line must start with #.
	if !strings.HasPrefix(line, atxHeadingChar) {
		return false
	}

	// Blank title is valid.
	if strings.Trim(line, atxHeadingChar+" ") == "" {
		return true
	}

	// The title line must have a space after #.
	return strings.HasPrefix(strings.TrimLeft(line, atxHeadingChar), " ")
}
