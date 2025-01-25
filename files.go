package fcqs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NotesFiles represents notes files.
type NotesFiles struct {
	Reader io.Reader
	Files  []*os.File
}

// Close closes all notes files.
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

	reader := io.MultiReader(readers...)

	return &NotesFiles{Reader: reader, Files: files}, nil
}

// notesFileNames returns filenames of notes.
func notesFileNames() ([]string, error) {
	f := os.Getenv("FCQS_NOTES_FILE")
	if f != "" {
		fileNames := strings.Split(f, string(os.PathListSeparator))
		return fileNames, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("user home directory: %w", err)
	}

	filenames := []string{filepath.Join(home, DefaultNotesFile)}
	return filenames, nil
}
