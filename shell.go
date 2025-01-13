package fcqs

import (
	_ "embed"
	"fmt"
	"io"
)

//go:embed shell.bash
var bashTemplate string

// WriteBashScript writes bash script to set up fcqs.
func WriteBashScript(w io.Writer) {
	fmt.Fprint(w, bashTemplate)
}
