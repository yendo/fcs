package value

import (
	"errors"
	"strings"
)

var ErrEmptyTitle = errors.New("title is empty")

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
		return nil, ErrEmptyTitle
	}

	return &Title{value: t}, nil
}
