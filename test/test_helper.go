package test

import "strings"

func GetExpectedTitles() string {
	titles := `title
long title one
title has regular expression meta chars $
contents have blank lines
same title
other heading level
title has trailing spaces
no contents
no contents2
no_space_title
fenced code block
URL
command-line
command-line with $
no blank line between title and contents
`

	return strings.Replace(titles, "trailing spaces", "trailing spaces  ", 1)
}
