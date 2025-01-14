package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yendo/fcqs/internal/value"
)

func TestFenceLineHasShellID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		line   string
		expect bool
	}{
		{name: "shell", line: "``` shell", expect: true},
		{name: "shell + string", line: "``` shell xxxx", expect: true},
		{name: "shell no space", line: "```shell", expect: true},
		{name: "shell more space", line: "```shell  ", expect: true},
		{name: "sh", line: "``` sh", expect: true},
		{name: "shell-script", line: "``` shell-script", expect: true},
		{name: "bash", line: "``` bash", expect: true},
		{name: "zsh", line: "``` zsh", expect: true},
		{name: "powershell", line: "``` powershell", expect: true},
		{name: "posh", line: "``` posh", expect: true},
		{name: "pwsh", line: "``` pwsh", expect: true},
		{name: "shellsession", line: "``` shellsession", expect: true},
		{name: "console", line: "``` console", expect: true},
		{name: "go", line: "``` go", expect: false},
		{name: "no identifier", line: "```", expect: false},
		{name: "other identifier", line: "``` other", expect: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fenceLine, ok := value.NewFenceLine(tc.line)

			require.True(t, ok)
			assert.Equal(t, tc.expect, fenceLine.HasShellID())
		})
	}
}

func TestFenceLineFuncs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		line   string
		expect bool
	}{
		{name: "normal", line: "```", expect: true},
		{name: "with trailing spaces", line: "```  ", expect: true},
		{name: "with identifier", line: "``` go", expect: true},
		{name: "with long identifier", line: "``` shell console", expect: true},
		{name: "no fence", line: "no fence", expect: false},
		{name: "no enough fence", line: "``", expect: false},
		{name: "head with spaces", line: " ```", expect: false},
		{name: "empty", line: "", expect: false},
	}

	t.Run("NewFenceLine", func(t *testing.T) {
		t.Parallel()

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				_, ok := value.NewFenceLine(tc.line)

				assert.Equal(t, tc.expect, ok)
			})
		}
	})

	t.Run("IsFenceLine", func(t *testing.T) {
		t.Parallel()

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				actual := value.IsFenceLine(tc.line)

				assert.Equal(t, tc.expect, actual)
			})
		}
	})
}
