package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func print_titles(buf io.Writer, fileName string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			fmt.Fprintln(buf, strings.TrimLeft(line, "# "))
		}
	}

	return nil
}

func print_contents(buf io.Writer, fileName string, title string) error {
	fp, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	flag := false
	prev := false

	r := regexp.MustCompile(fmt.Sprintf("^#* %s$", title))

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()

		if r.MatchString(line) {
			flag = true
		} else if flag && strings.HasPrefix(line, "#") {
			flag = false
		} else if flag && line == "" {
			prev = true
		}

		if flag && line != "" {
			if prev {
				fmt.Fprintln(buf, "")
				prev = false
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
		err = print_contents(os.Stdout, fileName, args[0])
	} else {
		err = print_titles(os.Stdout, fileName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
