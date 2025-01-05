package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/yendo/fcqs"
	"github.com/yendo/fcqs/internal/value"
)

var (
	version = "unknown"

	showVersion = flag.BoolP("version", "v", false, "output the version")
	showURL     = flag.BoolP("url", "u", false, "output the first URL from the note")
	showCmd     = flag.BoolP("command", "c", false, "output the first command from the note")
	showLoc     = flag.BoolP("location", "l", false, "output the note location")
	showBash    = flag.BoolP("bash", "", false, "output bash integration script")
	noTitle     = flag.BoolP("notitle", "t", false, "no title on output content")

	ErrInvalidNumberOfArgs = errors.New("invalid number of arguments")
)

func run(w io.Writer) error {
	var err error

	flag.Parse()
	args := flag.Args()

	if *showVersion {
		fmt.Fprintln(w, version)
		return nil
	}

	if *showBash {
		fcqs.WriteBashScript(w)
		return nil
	}

	notes, err := fcqs.OpenNotesFiles()
	if err != nil {
		return err
	}
	defer notes.Close()

	switch len(args) {
	case 0:
		if *showURL || *showCmd || *showLoc {
			return ErrInvalidNumberOfArgs
		}
		err = fcqs.WriteTitles(w, notes.Reader)
	case 1:
		title, tErr := value.NewTitle(args[0])
		if tErr != nil {
			// This error should be ignored to omit argument checking in shell scripts.
			return nil
		}

		switch {
		case *showURL:
			err = fcqs.WriteFirstURL(w, notes.Reader, title)
		case *showCmd:
			err = fcqs.WriteFirstCmdLineBlock(w, notes.Reader, title)
		case *showLoc:
			err = fcqs.WriteNoteLocation(w, notes.Files, title)
		default:
			err = fcqs.WriteContents(w, notes.Reader, title, *noTitle)
		}
	default:
		err = ErrInvalidNumberOfArgs
	}

	return err
}

func main() {
	exitCode := 0

	if err := run(os.Stdout); err != nil {
		exitCode = 1
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(exitCode)
}
