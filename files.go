package fcqs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type NotesFiles struct {
	Reader io.Reader
	Files  []*os.File
}

func (n NotesFiles) Close() {
	for _, f := range n.Files {
		f.Close()
	}
}

// NewNOtesFiles returns NotesFiles instance.
func OpenNotesFiles() (*NotesFiles, error) {
	fileName, err := notesFileNames()
	if err != nil {
		return nil, fmt.Errorf("notes file name: %w", err)
	}

	readers := make([]io.Reader, 0, len(fileName))
	files := make([]*os.File, 0, len(fileName))

	for _, v := range fileName {
		file, err := os.Open(v)
		if err != nil {
			return nil, fmt.Errorf("notes file: %w", err)
		}
		readers = append(readers, file)
		files = append(files, file)
	}

	file := io.MultiReader(readers...)

	return &NotesFiles{Reader: file, Files: files}, nil
}

// notesFileNames returns filenames of notes.
func notesFileNames() ([]string, error) {
	var fileNames []string

	f := os.Getenv("FCQS_NOTES_FILE")
	if f != "" {
		sep := ":"
		if runtime.GOOS == "windows" {
			sep = ";"
		}

		fileNames = strings.Split(f, sep)
		return fileNames, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("user home directory: %w", err)
	}

	fileNames = append(fileNames, filepath.Join(home, DefaultNotesFile))
	return fileNames, nil
}
