// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := LoadTheme("../../testdata/src", "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 2).
		Equal(theme.ID, "default")

	theme, err = LoadTheme("../../testdata/src", "not-exists")
	a.ErrorIs(err, os.ErrNotExist).Nil(theme)
}

func TestTheme_sanitize(t *testing.T) {
	a := assert.New(t)
	theme := &Theme{Index: "index.xsl", Tags: "tags.xsl", Tag: "tag.xsl", Templates: []string{"post.xsl"}}
	a.NotError(theme.sanitize("../../testdata/src/themes/default", "default"))
	a.Equal(theme.ID, "default").
		Empty(theme.Description).
		Equal(theme.Templates, []string{"post.xsl"}).
		Equal(theme.Index, "index.xsl")

	theme = &Theme{Index: "index.xsl", Tags: "tags.xsl", Tag: "tag.xsl", Templates: []string{"style.xsl", "not-exists"}}
	err := theme.sanitize("../../testdata/src/themes/default", "default")
	a.Error(err).Equal(err.Field, "templates.not-exists")

	theme = &Theme{Index: "index.xsl", Tags: "tags.xsl", Tag: "tag.xsl", Templates: []string{"style.xsl"}, Screenshots: []string{"not-exists"}}
	err = theme.sanitize("../../testdata/src/themes/default", "default")
	a.Error(err).Equal(err.Field, "screenshots[0]")
}
