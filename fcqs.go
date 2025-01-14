package fcqs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/yendo/fcqs/internal/value"
	"mvdan.cc/xurls/v2"
)

const (
	DefaultNotesFile = "fcnotes.md"

	shellPrompt = "$"

	// State of text line.
	normal = iota
	fenced
	scoped
	scopedFenced
)

// WriteTitles writes the titles of all notes.
func WriteTitles(w io.Writer, r io.Reader) error {
	var allTitles []value.Title
	var title value.Title

	scanner := bufio.NewScanner(r)
	state := normal

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		switch state {
		case normal:
			if titleLine, ok := value.NewTitleLine(line); ok {
				if titleLine.HasValidTitle() {
					title = titleLine.Title()
					continue
				}
			}

			if !slices.Contains(allTitles, title) {
				fmt.Fprintln(w, title)
				allTitles = append(allTitles, title)
			}

			if value.IsFenceLine(line) {
				state = fenced
			}

		case fenced:
			if value.IsFenceLine(line) {
				state = normal
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
	f := newFilter(w, isNoTitle)
	defer f.Close()

	state := normal

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case normal:
			if value.IsFenceLine(line) {
				state = fenced
			}

			if titleLine, ok := value.NewTitleLine(line); ok {
				if titleLine.EqualTitle(title) {
					state = scoped
					fmt.Fprint(f, line)
				}
			}

		case fenced:
			if value.IsFenceLine(line) {
				state = normal
			}

		case scoped:
			if titleLine, ok := value.NewTitleLine(line); ok {
				if !titleLine.EqualTitle(title) {
					state = normal
					break
				}
			}

			if value.IsFenceLine(line) {
				state = scopedFenced
			}

			fmt.Fprint(f, line)

		case scopedFenced:
			if value.IsFenceLine(line) {
				state = scoped
			}
			fmt.Fprint(f, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("seek contents: %w", err)
	}

	return nil
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
	state := normal

	var buf bytes.Buffer
	if err := WriteContents(&buf, r, title, false); err != nil {
		return err
	}
	scanner := bufio.NewScanner(&buf)

	for scanner.Scan() {
		line := scanner.Text()
		fenceLine, isFenceLine := value.NewFenceLine(line)

		switch state {
		case normal:
			if isFenceLine && fenceLine.HasShellID() {
				state = fenced
			}

		case fenced:
			if isFenceLine {
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

// WriteNoteLocation writes the file name and line number of the note.
func WriteNoteLocation(w io.Writer, files []*os.File, title *value.Title) error {
	for _, file := range files {
		c := 0
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			c++
			line := scanner.Text()

			if titleLine, ok := value.NewTitleLine(line); ok {
				if titleLine.EqualTitle(title) {
					fmt.Fprintf(w, "%q %d\n", file.Name(), c)
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("seek note location: %w", err)
		}
	}
	return nil
}
