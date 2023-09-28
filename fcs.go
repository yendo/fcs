package fcs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"mvdan.cc/xurls/v2"
)

const DefaultNotesFile = "fcnotes.md"

// WriteTitles writes the titles of all notes.
func WriteTitles(w io.Writer, r io.Reader) {
	var allTitles []string

	scanner := bufio.NewScanner(r)
	isFenced := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") && !isFenced {
			// skip titles without a space after the `#`
			if !strings.HasPrefix(strings.TrimLeft(line, "#"), " ") {
				continue
			}

			title := strings.Trim(line, "# ")
			if title == "" {
				continue
			}

			if !slices.Contains(allTitles, title) {
				fmt.Fprintln(w, title)
				allTitles = append(allTitles, title)
			}
		}

		if strings.HasPrefix(line, "```") {
			isFenced = !isFenced
		}
	}
}

// WriteContents writes the contents of the note.
func WriteContents(w io.Writer, r io.Reader, title string) {
	title = strings.Trim(title, " ")

	isScope := false
	isFenced := false
	isBlank := false
	re := getNoteTitleRegexp(title)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if !isFenced && re.MatchString(line) {
			isScope = true
		} else if isScope {
			switch {
			case !isFenced && isContentsEnd(line):
				isScope = false
			case strings.HasPrefix(line, "```"):
				isFenced = !isFenced
			case line == "":
				isBlank = true
			}
		}

		if isScope && line != "" {
			if isBlank {
				isBlank = false

				fmt.Fprintln(w, "")
			}

			fmt.Fprintln(w, line)
		}
	}
}

// isContentsEnd returns if the line is the end of the contents.
func isContentsEnd(line string) bool {
	// Title must start with #.
	if !strings.HasPrefix(line, "#") {
		return false
	}

	// Title may be blank.
	if strings.Trim(line, "# ") == "" {
		return true
	}

	// Title must have a space after #
	return strings.HasPrefix(strings.TrimLeft(line, "#"), " ")
}

// PrintsFirstURL writes the first URL in the contents of the note.
func WriteFirstURL(w io.Writer, r io.Reader, title string) {
	var buf bytes.Buffer

	WriteContents(&buf, r, title)

	rxStrict := xurls.Strict()
	url := rxStrict.FindString(buf.String())

	if url != "" {
		fmt.Fprintln(w, url)
	}
}

// WriteFirstCmdLine writes the first command-line in the contents of the note.
func WriteFirstCmdLine(w io.Writer, r io.Reader, title string) {
	var buf bytes.Buffer

	isFenced := false

	WriteContents(&buf, r, title)
	scanner := bufio.NewScanner(&buf)

	for scanner.Scan() {
		line := scanner.Text()

		if isShellCodeBlockBegin(line) {
			isFenced = true

			continue
		} else if strings.HasPrefix(line, "```") && isFenced {
			break
		}

		if isFenced {
			fmt.Fprintln(w, strings.TrimLeft(line, "$ "))
		}
	}
}

var reShellCodeBlock = regexp.MustCompile("^```+\\s*(\\S+).*$")

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
func WriteNoteLocation(w io.Writer, file *os.File, title string) {
	c := 0
	scanner := bufio.NewScanner(file)
	r := getNoteTitleRegexp(title)

	for scanner.Scan() {
		c++
		line := scanner.Text()

		if r.MatchString(line) {
			fmt.Fprintf(w, "%q %d\n", file.Name(), c)

			break
		}
	}
}

// GetNotesFileName returns the filename of the notes.
func GetNotesFileName() (string, error) {
	fileName := os.Getenv("FCS_NOTES_FILE")
	if fileName != "" {
		return fileName, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot access user home directory: %w", err)
	}

	fileName = filepath.Join(home, DefaultNotesFile)

	return fileName, nil
}

// getNoteTitleRegexp returns a regular expression to search for the title of the note.
func getNoteTitleRegexp(title string) *regexp.Regexp {
	if title == "" {
		return regexp.MustCompile(fmt.Sprintf("^#+\\s*%s\\s*$", regexp.QuoteMeta(title)))
	} else {
		return regexp.MustCompile(fmt.Sprintf("^#+\\s+%s\\s*$", regexp.QuoteMeta(title)))
	}
}
