package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"mvdan.cc/xurls/v2"
)

var (
	version = "unknown"

	showVersion = flag.Bool("v", false, "output the version")
	showURL     = flag.Bool("u", false, "output the first URL from the note")
	showCmd     = flag.Bool("c", false, "output the first command from the note")
	showLoc     = flag.Bool("l", false, "output the note location")

	ErrInvalidNumberOfArgs = errors.New("invalid number of arguments")
)

func printTitles(buf io.Writer, fd io.Reader) {
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

func printContents(buf io.Writer, fd io.Reader, title string) {
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

func printFirstURL(buf io.Writer, fd io.Reader, title string) {
	var b bytes.Buffer

	printContents(&b, fd, title)

	rxStrict := xurls.Strict()
	url := rxStrict.FindString(b.String())

	if url != "" {
		fmt.Fprintln(buf, url)
	}
}

func printFirstCmdLine(buf io.Writer, fd io.Reader, title string) {
	var b bytes.Buffer

	isFenced := false

	printContents(&b, fd, title)
	scanner := bufio.NewScanner(&b)

	for scanner.Scan() {
		line := scanner.Text()

		if isShellCodeBlockBegin(line) {
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

func printLineNumber(buf io.Writer, fd *os.File, title string) {
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

func getNotesFile() (string, error) {
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

func run(buf io.Writer) error {
	flag.Parse()
	args := flag.Args()

	if *showVersion {
		fmt.Fprintln(buf, version)

		return nil
	}

	fileName, err := getNotesFile()
	if err != nil {
		return err
	}

	fd, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("cannot access notes file: %w", err)
	}
	defer fd.Close()

	if *showURL || *showCmd || *showLoc {
		if len(args) != 1 {
			return ErrInvalidNumberOfArgs
		}

		switch {
		case *showURL:
			printFirstURL(buf, fd, args[0])
		case *showCmd:
			printFirstCmdLine(buf, fd, args[0])
		case *showLoc:
			printLineNumber(buf, fd, args[0])
		}

		return nil
	}

	switch len(args) {
	case 0:
		printTitles(buf, fd)
	case 1:
		printContents(buf, fd, args[0])
	default:
		return ErrInvalidNumberOfArgs
	}

	return nil
}

func main() {
	exitCode := 0

	if err := run(os.Stdout); err != nil {
		exitCode = 1

		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(exitCode)
}
