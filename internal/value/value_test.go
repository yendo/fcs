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

func TestNewTitleLine(t *testing.T) {
	t.Parallel()

	t.Run("returns True", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name       string
			titleLine  string
			validTitle bool
		}{
			{name: "normal title", titleLine: "# title string", validTitle: true},
			{name: "blank title", titleLine: "#   ", validTitle: false},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				titleLine, ok := value.NewTitleLine(tc.titleLine)

				require.True(t, ok)
				assert.Equal(t, tc.validTitle, titleLine.HasValidTitle())
			})
		}
	})

	t.Run("returns False", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			titleLine string
		}{
			{name: "no title prefix", titleLine: "no title prefix"},
			{name: "vacant line", titleLine: ""},
			{name: "only spaces", titleLine: "  "},
			{name: "no space", titleLine: "#no_space"},
			{name: "fenced chars", titleLine: "```"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				titleLine, ok := value.NewTitleLine(tc.titleLine)

				require.False(t, ok)
				assert.Nil(t, titleLine)
			})
		}
	})
}

func TestTitleLineTitle(t *testing.T) {
	t.Parallel()

	titleLine, ok := value.NewTitleLine("# sample title")
	assert.True(t, ok)

	title := titleLine.Title()
	assert.Equal(t, "sample title", title.String())
}

func TestTitleLineEqualTitle(t *testing.T) {
	t.Parallel()

	titleLine, ok := value.NewTitleLine("# sample title")
	assert.True(t, ok)

	title, err := value.NewTitle("sample title")
	require.NoError(t, err)

	otherTitle, err := value.NewTitle("other title")
	require.NoError(t, err)

	assert.True(t, titleLine.EqualTitle(title))
	assert.False(t, titleLine.EqualTitle(otherTitle))
}
