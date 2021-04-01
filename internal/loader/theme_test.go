// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/testdata"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := LoadTheme(testdata.Source, "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 2).
		Equal(theme.ID, "default").
		Equal(3, len(theme.Highlights))

	theme, err = LoadTheme(testdata.Source, "not-exists")
	a.ErrorIs(err, fs.ErrNotExist).Nil(theme)
}

func TestTheme_sanitize(t *testing.T) {
	a := assert.New(t)

	theme := &Theme{Templates: []string{"post"}}
	a.NotError(theme.sanitize(testdata.Source, "themes/default", "default"))
	a.Equal(theme.ID, "default").
		Empty(theme.Description).
		Equal(theme.Templates, []string{"post"})

	theme = &Theme{Templates: []string{"style.xsl"}, Screenshots: []string{"not-exists"}}
	err := theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "screenshots[0]")
}

func TestHighlight_sanitize(t *testing.T) {
	a := assert.New(t)

	h := &Highlight{}
	err := h.sanitize()
	a.Equal(err.Field, "name")

	h = &Highlight{Name: "not-exists"}
	err = h.sanitize()
	a.Equal(err.Field, "name")

	h = &Highlight{Name: "solarized-dark256"}
	a.NotError(h.sanitize())
}
