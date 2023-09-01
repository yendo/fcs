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

func printTitles(buf io.Writer, fileName string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	var allTitles []string

	scanner := bufio.NewScanner(fp)
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

	return nil
}

func printContents(buf io.Writer, fileName string, title string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	isScope := false
	isBlank := false

	r := regexp.MustCompile(fmt.Sprintf("^#* %s$", title))

	scanner := bufio.NewScanner(fp)
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

	return nil
}

func main() {
	flag.Parse()
	args := flag.Args()
	var err error

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fileName := filepath.Join(home, "fcmemo.md")
	if len(args) == 1 {
		err = printContents(os.Stdout, fileName, args[0])
	} else {
		err = printTitles(os.Stdout, fileName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
