// SPDX-License-Identifier: MIT

package data

import (
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Sitemap 的相关配置项
type Sitemap struct {
	*loader.Sitemap
	Permalink    string
	XSLPermalink string
	Path         string
}

func newSitemap(conf *loader.Config, theme *loader.Theme) *Sitemap {
	p := buildPath(vars.SitemapXML)
	sm := &Sitemap{
		Sitemap:   conf.Sitemap,
		Permalink: BuildURL(conf.URL, p),
		Path:      p,
	}

	if theme.Sitemap != "" {
		sm.XSLPermalink = buildThemeURL(conf.URL, conf.Theme, theme.Sitemap)
	}

	return sm
}
