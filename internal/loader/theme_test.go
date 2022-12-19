// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2/internal/testdata"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t, false)

	theme, err := LoadTheme(testdata.Source, "default")
	a.NotError(err).NotNil(theme)
	a.Equal(len(theme.Authors), 1).
		Equal(theme.ID, "default").
		Equal(3, len(theme.Highlights))

	theme, err = LoadTheme(testdata.Source, "not-exists")
	a.ErrorIs(err, fs.ErrNotExist).Nil(theme)
}

func TestTheme_sanitize(t *testing.T) {
	a := assert.New(t, false)

	theme := &Theme{Templates: []string{"post"}}
	a.NotError(theme.sanitize(testdata.Source, "themes/default", "default"))
	a.Equal(theme.ID, "default").
		Empty(theme.Description).
		Equal(theme.Templates, []string{"post"})

	// rss 不存在
	theme = &Theme{Templates: []string{"style.xsl"}, RSS: "not-exists"}
	err := theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "rss")

	// atom 不存在
	theme = &Theme{Templates: []string{"style.xsl"}, Atom: "not-exists"}
	err = theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "atom")

	// sitemap 不存在
	theme = &Theme{Templates: []string{"style.xsl"}, Sitemap: "not-exists"}
	err = theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "sitemap")

	// screenshots 不存在
	theme = &Theme{Templates: []string{"style.xsl"}, Screenshots: []string{"not-exists"}}
	err = theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "screenshots[0]")

	// highlight.name 为空
	theme = &Theme{Templates: []string{"style.xsl"}, Highlights: []*Highlight{{}}}
	err = theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "highlight[0].name")

	// 两个空的 Highlight.Media
	theme = &Theme{Templates: []string{"style.xsl"}, Highlights: []*Highlight{
		{Name: "bw"},
		{Name: "algol"},
	}}
	err = theme.sanitize(testdata.Source, "themes/default", "default")
	a.Error(err).Equal(err.Field, "highlight[1].media") // 第二个元素的 media 不能为空
}

func TestHighlight_sanitize(t *testing.T) {
	a := assert.New(t, false)

	h := &Highlight{}
	err := h.sanitize()
	a.Equal(err.Field, "name")

	h = &Highlight{Name: "not-exists"}
	err = h.sanitize()
	a.Equal(err.Field, "name")

	h = &Highlight{Name: "solarized-dark256"}
	a.NotError(h.sanitize())
}
