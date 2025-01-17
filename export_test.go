package fcqs

import (
	"bufio"
	"io"
	"testing"
	"testing/iotest"
)

var ExportNewFilter = newFilter

// SetNewScannerMock sets error mock for bufio.NewScanner.
func SetNewScannerMock(t *testing.T, e error) {
	t.Helper()

	var tmp, s func(r io.Reader) *bufio.Scanner

	s = func(_ io.Reader) *bufio.Scanner {
		return bufio.NewScanner(iotest.ErrReader(e))
	}
	tmp, newScanner = newScanner, s

	t.Cleanup(func() {
		newScanner = tmp
	})
}
