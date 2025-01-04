package value

import (
	"fmt"
	"strings"
)

type Title struct {
	value string
}

func (t Title) String() string {
	return t.value
}

func NewTitle(t string) (*Title, error) {
	t = strings.Trim(t, " ")
	if t == "" {
		return nil, fmt.Errorf("title is empty")
	}

	return &Title{value: t}, nil
}
