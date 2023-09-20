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

func PrintTitles(buf io.Writer, fd io.Reader) {
	var allTitles []string

	scanner := bufio.NewScanner(fd)
	isFenced := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") && !isFenced {
			title := strings.TrimLeft(line, "# ")
			if title == "" {
				continue
			}

			if !slices.Contains(allTitles, title) {
				fmt.Fprintln(buf, title)
				allTitles = append(allTitles, title)
			}
		}

		if strings.HasPrefix(line, "```") {
			isFenced = !isFenced
		}
	}
}

func PrintContents(buf io.Writer, fd io.Reader, title string) {
	isScope := false
	isFenced := false
	isBlank := false
	r := getNoteTitleRegexp(title)

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()

		if r.MatchString(line) && !isFenced {
			isScope = true
		} else if isScope {
			switch {
			case strings.HasPrefix(line, "#") && !isFenced:
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

				fmt.Fprintln(buf, "")
			}

			fmt.Fprintln(buf, line)
		}
	}
}

func PrintFirstURL(buf io.Writer, fd io.Reader, title string) {
	var b bytes.Buffer

	PrintContents(&b, fd, title)

	rxStrict := xurls.Strict()
	url := rxStrict.FindString(b.String())

	if url != "" {
		fmt.Fprintln(buf, url)
	}
}

func PrintFirstCmdLine(buf io.Writer, fd io.Reader, title string) {
	var b bytes.Buffer

	isFenced := false

	PrintContents(&b, fd, title)
	scanner := bufio.NewScanner(&b)

	for scanner.Scan() {
		line := scanner.Text()

		if IsShellCodeBlockBegin(line) {
			isFenced = true

			continue
		} else if strings.HasPrefix(line, "```") && isFenced {
			break
		}

		if isFenced {
			fmt.Fprintln(buf, strings.TrimLeft(line, "$ "))
		}
	}
}

var reShellCodeBlock = regexp.MustCompile("^```\\s*(\\S+).*$")

func IsShellCodeBlockBegin(line string) bool {
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

func PrintLineNumber(buf io.Writer, fd *os.File, title string) {
	c := 0
	scanner := bufio.NewScanner(fd)
	r := getNoteTitleRegexp(title)

	for scanner.Scan() {
		c++
		line := scanner.Text()

		if r.MatchString(line) {
			fmt.Fprintf(buf, "%q %d\n", fd.Name(), c)

			break
		}
	}
}

func GetNotesFile() (string, error) {
	fileName := os.Getenv("FCS_NOTES_FILE")
	if fileName != "" {
		return fileName, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot access user home directory: %w", err)
	}

	fileName = filepath.Join(home, "fcnotes.md")

	return fileName, nil
}

func getNoteTitleRegexp(title string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^#* %s$", regexp.QuoteMeta(title)))
}
