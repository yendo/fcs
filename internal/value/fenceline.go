package value

import (
	"slices"
	"strings"
)

const codeFence = "```"

var shellList = []string{
	"shell", "sh", "shell-script", "bash", "zsh",
	"powershell", "posh", "pwsh",
	"shellsession", "console",
}

// FenceLine represents a fence text line.
type FenceLine struct {
	line string
}

// HasShellID reports whether the fence line has shell identifier.
func (fl FenceLine) HasShellID() bool {
	trimmedLine := strings.Trim(fl.line, "` ")
	if trimmedLine == "" {
		return false
	}

	id := strings.Split(trimmedLine, " ")
	return slices.Contains(shellList, id[0])
}

// NewFenceLine returns Fence line.
func NewFenceLine(line string) (*FenceLine, bool) {
	if !IsFenceLine(line) {
		return nil, false
	}

	return &FenceLine{line: line}, true
}

// IsFenceLine reports whether the line is fence line.
func IsFenceLine(line string) bool {
	return strings.HasPrefix(line, codeFence)
}
