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

var (
	version     string = "unknown"
	showVersion *bool  = flag.Bool("v", false, "Show version")
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

	r := regexp.MustCompile(fmt.Sprintf("^#* %s$", title))

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
