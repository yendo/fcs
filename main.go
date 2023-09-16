package main

import (
	"bufio"
	"bytes"
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
	version     = "unknown"
	showVersion = flag.Bool("v", false, "output version")
	showURL     = flag.Bool("u", false, "output first URL from a note")
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

	r := regexp.MustCompile(fmt.Sprintf("^#* %s$", regexp.QuoteMeta(title)))

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()

		if r.MatchString(line) && !isFenced {
			isScope = true
		} else if isScope {
			if strings.HasPrefix(line, "#") && !isFenced {
				isScope = false
			} else if strings.HasPrefix(line, "```") {
				isFenced = !isFenced
			} else if line == "" {
				isBlank = true
			}
		}

		if isScope && line != "" {
			if isBlank {
				fmt.Fprintln(buf, "")
				isBlank = false
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

func getNotesFile() (string, error) {
	fileName := os.Getenv("FCS_NOTES_FILE")
	if fileName != "" {
		return fileName, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fileName = filepath.Join(home, "fcnotes.md")
	return fileName, nil
}

func run(buf io.Writer) error {
	flag.Parse()
	args := flag.Args()
	var err error

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
		return err
	}
	defer fd.Close()

	if *showURL {
		if len(args) != 1 {
			return fmt.Errorf("invalid number of arguments")
		}
		printFirstURL(buf, fd, args[0])
	}

	switch len(args) {
	case 0:
		printTitles(buf, fd)
	case 1:
		printContents(buf, fd, args[0])
	default:
		return fmt.Errorf("invalid number of arguments")
	}

	return nil
}

func main() {
	exitCode := 0

	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
