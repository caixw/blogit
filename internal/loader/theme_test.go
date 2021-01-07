// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := LoadTheme("../testdata", "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 2).
		Equal(theme.ID, "default")

	theme, err = LoadTheme("../testdata", "not-exists")
	a.ErrorIs(err, os.ErrNotExist).Nil(theme)
}

func TestTheme_sanitize(t *testing.T) {
	a := assert.New(t)
	theme := &Theme{}
	a.NotError(theme.sanitize("../testdata", "def"))
	a.Equal(theme.ID, "def").
		Empty(theme.Description).
		Equal(theme.Templates, []string{"post.xsl"})

	theme = &Theme{Templates: []string{"style.xsl"}}
	a.NotError(theme.sanitize("../testdata/themes/default", "default"))

	theme = &Theme{Templates: []string{"style.xsl", "not-exists"}}
	err := theme.sanitize("../testdata/themes/default", "default")
	a.Error(err).Equal(err.Field, "templates.not-exists")

	theme = &Theme{Screenshots: []string{"not-exists"}}
	err = theme.sanitize("../testdata/themes/default", "default")
	a.Error(err).Equal(err.Field, "screenshots[0]")
}
