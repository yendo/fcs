package fcqs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcqs"
	"github.com/yendo/fcqs/test"
)

func TestNewNotesFiles(t *testing.T) {
	t.Run("failed to access user home directory", func(t *testing.T) {
		t.Setenv("FCQS_NOTES_FILE", "")
		t.Setenv("HOME", "")

		notes, err := fcqs.OpenNotesFiles()

		require.Error(t, err)
		require.EqualError(t, err, "notes file name: user home directory: $HOME is not defined")
		assert.Nil(t, notes)
	})

	t.Run("failed to access notes file", func(t *testing.T) {
		expectedFileNames := test.MultiFiles(test.NotesFile, "invalid_file")
		t.Setenv("FCQS_NOTES_FILE", expectedFileNames)

		notes, err := fcqs.OpenNotesFiles()

		require.Error(t, err)
		require.EqualError(t, err, "notes file: open invalid_file: no such file or directory")
		assert.Nil(t, notes)
	})

	t.Run("set a file from environment variable", func(t *testing.T) {
		expectedFileName := test.NotesFile
		t.Setenv("FCQS_NOTES_FILE", expectedFileName)

		notes, err := fcqs.OpenNotesFiles()
		notes.Close()

		require.NoError(t, err)
		assert.NotNil(t, notes.Reader)
		assert.Equal(t, expectedFileName, notes.Files[0].Name())
	})

	t.Run("set files from environment variable", func(t *testing.T) {
		expectedFileNames := []string{test.NotesFile, test.LocationFile, test.LocationExtraFile}
		t.Setenv("FCQS_NOTES_FILE", test.MultiFiles(expectedFileNames...))

		notes, err := fcqs.OpenNotesFiles()
		notes.Close()

		require.NoError(t, err)
		assert.NotNil(t, notes.Reader)
		assert.Equal(t, expectedFileNames[0], notes.Files[0].Name())
		assert.Equal(t, expectedFileNames[1], notes.Files[1].Name())
		assert.Equal(t, expectedFileNames[2], notes.Files[2].Name())
	})

	t.Run("default filename", func(t *testing.T) {
		// Arrange
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)
		t.Setenv("FCQS_NOTES_FILE", "")

		file := filepath.Join(tempHome, fcqs.DefaultNotesFile)
		err := os.WriteFile(file, []byte("# title\ncontents\n"), 0o600)
		require.NoError(t, err)

		home, err := os.UserHomeDir()
		require.NoError(t, err)
		require.Equal(t, file, filepath.Join(home, fcqs.DefaultNotesFile))

		// Act
		notes, err := fcqs.OpenNotesFiles()
		notes.Close()

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, notes.Reader)
		assert.Equal(t, filepath.Join(home, fcqs.DefaultNotesFile), notes.Files[0].Name())
	})
}
