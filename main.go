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

func run(buf io.Writer) error {
	flag.Parse()
	args := flag.Args()
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fileName := filepath.Join(home, "fcnotes.md")
	fd, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	if len(args) == 1 {
		printContents(buf, fd, args[0])
	} else {
		printTitles(buf, fd)
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
