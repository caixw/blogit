// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)
	fs := os.DirFS("../../testdata/src")

	theme, err := LoadTheme(fs, "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 2).
		Equal(theme.ID, "default")

	theme, err = LoadTheme(fs, "not-exists")
	a.ErrorIs(err, os.ErrNotExist).Nil(theme)
}

func TestTheme_sanitize(t *testing.T) {
	a := assert.New(t)
	fs := os.DirFS("../../testdata/src")

	theme := &Theme{Templates: []string{"post"}}
	a.NotError(theme.sanitize(fs, "themes/default", "default"))
	a.Equal(theme.ID, "default").
		Empty(theme.Description).
		Equal(theme.Templates, []string{"post"})

	theme = &Theme{Templates: []string{"style.xsl"}, Screenshots: []string{"not-exists"}}
	err := theme.sanitize(fs, "themes/default", "default")
	a.Error(err).Equal(err.Field, "screenshots[0]")
}
