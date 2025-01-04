package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yendo/fcqs/internal/value"
)

func TestNewTitle(t *testing.T) {
	t.Parallel()

	t.Run("success cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			title string
		}{
			{name: "trimmed title", title: "title string"},
			{name: "un-trimmed title", title: " title string  "},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				title, err := value.NewTitle(tc.title)

				assert.NoError(t, err)
				assert.Equal(t, "title string", title.String())
			})
		}
	})

	t.Run("fail cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			title string
		}{
			{name: "empty title", title: ""},
			{name: "white spaces title", title: "  "},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				title, err := value.NewTitle(tc.title)

				assert.Error(t, err)
				assert.Nil(t, title)
			})
		}
	})
}
