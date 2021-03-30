// SPDX-License-Identifier: MIT

package data

import (
	"path"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Theme 主题描述
type Theme struct {
	ID          string
	URL         string
	Description string
	Authors     []*loader.Author
}

// Highlight 高亮主题
type Highlight struct {
	Name  string // 主题名称
	Path  string
	URL   string
	Media string
}

func newTheme(t *loader.Theme) *Theme {
	return &Theme{
		ID:          t.ID,
		URL:         t.URL,
		Description: t.Description,
		Authors:     t.Authors,
	}
}

func newHighlights(conf *loader.Config, t *loader.Theme) []*Highlight {
	hs := make([]*Highlight, 0, len(t.Highlights))
	for _, h := range t.Highlights {
		p := path.Join(vars.ThemesDir, t.ID, h.Name+".css")
		hs = append(hs, &Highlight{
			Name:  h.Name,
			Path:  p,
			URL:   BuildURL(conf.URL, p),
			Media: h.Media,
		})
	}

	return hs
}
