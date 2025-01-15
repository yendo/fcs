package value_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

				require.NoError(t, err)
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

				require.Error(t, err)
				assert.Nil(t, title)
			})
		}
	})
}

func TestTitleEquals(t *testing.T) {
	t.Parallel()

	title1, err := value.NewTitle("sample title")
	require.NoError(t, err)

	title2, err := value.NewTitle("sample title")
	require.NoError(t, err)

	otherTitle, err := value.NewTitle("other title")
	require.NoError(t, err)

	assert.True(t, title1.Equals(title2))
	assert.False(t, title1.Equals(otherTitle))
}
