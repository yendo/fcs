package fcqs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/yendo/fcqs/internal/value"
	"mvdan.cc/xurls/v2"
)

const (
	DefaultNotesFile = "fcnotes.md"

	codeFence      = "```"
	atxHeadingChar = "#"
	shellPrompt    = "$"

	// State of text line.
	Normal = iota
	Fenced
	Scoped
	ScopedFenced
)

// WriteTitles writes the titles of all notes.
func WriteTitles(w io.Writer, r io.Reader) error {
	var allTitles []string
	var title string

	scanner := bufio.NewScanner(r)
	state := Normal

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		switch state {
		case Normal:
			if isTitleLineWithString(line) {
				title = trimmedTitle(line)
				continue
			}

			if !slices.Contains(allTitles, title) {
				fmt.Fprintln(w, title)
				allTitles = append(allTitles, title)
			}

			if strings.HasPrefix(line, codeFence) {
				state = Fenced
			}

		case Fenced:
			if strings.HasPrefix(line, codeFence) {
				state = Normal
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("seek titles: %w", err)
	}

	return nil
}

// WriteContents writes the contents of the note.
func WriteContents(w io.Writer, r io.Reader, title *value.Title, isNoTitle bool) error {
	state := Normal

	f := NewFilter(w, isNoTitle)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case Normal:
			if strings.HasPrefix(line, codeFence) {
				state = Fenced
			}

			if isSearchedTitleLine(line, title) {
				state = Scoped
				f.Write(line)
			}

		case Fenced:
			if strings.HasPrefix(line, codeFence) {
				state = Normal
			}

		case Scoped:
			if isTitleLine(line) && !isSearchedTitleLine(line, title) {
				state = Normal
				break
			}

			if strings.HasPrefix(line, codeFence) {
				state = ScopedFenced
			}

			f.Write(line)

		case ScopedFenced:
			if strings.HasPrefix(line, codeFence) {
				state = Scoped
			}
			f.Write(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("seek contents: %w", err)
	}

	f.Write("")
	return nil
}

// isTitleLine returns if the line is title line.
func isTitleLine(line string) bool {
	// Title line must start with #.
	if !strings.HasPrefix(line, atxHeadingChar) {
		return false
	}

	// Blank title is valid.
	if trimmedTitle(line) == "" {
		return true
	}

	// The title line must have a space after #.
	return strings.HasPrefix(strings.TrimLeft(line, atxHeadingChar), " ")
}

// isTitleLineWithString returns if the line is title line with string.
func isTitleLineWithString(line string) bool {
	// Title line with string should be title line.
	if !isTitleLine(line) {
		return false
	}

	return trimmedTitle(line) != ""
}

// isSearchedTitleLine returns if the line is the searched title line.
func isSearchedTitleLine(line string, title *value.Title) bool {
	// Searched title line should be title line.
	if !isTitleLine(line) {
		return false
	}

	// When the trimmed line and title match, the content starts.
	return trimmedTitle(line) == title.String()
}

func trimmedTitle(line string) string {
	return strings.Trim(line, atxHeadingChar+" ")
}

// WriteFirstURL writes the first URL in the contents of the note.
func WriteFirstURL(w io.Writer, r io.Reader, title *value.Title) error {
	var buf bytes.Buffer
	if err := WriteContents(&buf, r, title, false); err != nil {
		return err
	}

	rxStrict := xurls.Strict()
	url := rxStrict.FindString(buf.String())

	if url != "" {
		fmt.Fprintln(w, url)
	}

	return nil
}

// WriteFirstCmdLineBlock writes the first command-line block in the contents of the note.
func WriteFirstCmdLineBlock(w io.Writer, r io.Reader, title *value.Title) error {
	state := Normal

	var buf bytes.Buffer
	if err := WriteContents(&buf, r, title, false); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&buf)

	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case Normal:
			if isShellCodeBlockBegin(line) {
				state = Fenced
			}

		case Fenced:
			if strings.HasPrefix(line, codeFence) {
				break
			}
			fmt.Fprintln(w, strings.TrimLeft(line, shellPrompt+" "))
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("seek command line block: %w", err)
	}

	return nil
}

var reShellCodeBlock = regexp.MustCompile(fmt.Sprintf("^%s+\\s*(\\S+).*$", codeFence))

// isShellCodeBlockBegin determines if the line is the beginning of a shell code block.
func isShellCodeBlockBegin(line string) bool {
	shellList := []string{
		"shell", "sh", "shell-script", "bash", "zsh",
		"powershell", "posh", "pwsh",
		"shellsession", "console",
	}

	match := reShellCodeBlock.FindStringSubmatch(line)
	if len(match) == 0 {
		return false
	}

	return slices.Contains(shellList, match[1])
}

// WriteNoteLocation writes the file name and line number of the note.
func WriteNoteLocation(w io.Writer, files []*os.File, title *value.Title) error {
	for _, file := range files {
		c := 0
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			c++
			line := scanner.Text()

			if isSearchedTitleLine(line, title) {
				fmt.Fprintf(w, "%q %d\n", file.Name(), c)
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("seek note location: %w", err)
		}
	}
	return nil
}
