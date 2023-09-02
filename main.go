package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const version string = "0.1.0"

var showVersion = flag.Bool("v", false, "Show version")

func printTitles(buf io.Writer, fd io.Reader) {
	var allTitles []string

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			title := strings.TrimLeft(line, "# ")
			if title == "" {
				continue
			}

			if !slices.Contains(allTitles, title) {
				fmt.Fprintln(buf, title)
				allTitles = append(allTitles, title)
			}
		}
	}
}

func printContents(buf io.Writer, fd io.Reader, title string) {
	isScope := false
	isBlank := false

	r := regexp.MustCompile(fmt.Sprintf("^#* %s$", title))

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()

		if r.MatchString(line) {
			isScope = true
		} else if isScope {
			if strings.HasPrefix(line, "#") {
				isScope = false
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

	err := run(os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}

	os.Exit(exitCode)
}
