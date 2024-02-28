// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package data

import (
	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
)

// Sitemap 的相关配置项
type Sitemap struct {
	*loader.Sitemap
	Permalink    string
	XSLPermalink string
	Path         string
}

func newSitemap(conf *loader.Config, theme *loader.Theme) *Sitemap {
	sm := &Sitemap{
		Sitemap:   conf.Sitemap,
		Permalink: BuildURL(conf.URL, vars.SitemapXML),
		Path:      vars.SitemapXML,
	}

	if theme.Sitemap != "" {
		sm.XSLPermalink = buildThemeURL(conf.URL, conf.Theme, theme.Sitemap)
	}

	return sm
}
