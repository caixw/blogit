// SPDX-License-Identifier: MIT

// Package data 对加载的数据进行二次加工
package data

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

type (
	// Data 处理后的数据
	Data struct {
		URL         string
		Title       string
		Subtitle    string
		TitleSuffix string // 每篇文章标题的后缀
		Icon        *loader.Icon
		Language    string
		Authors     []*loader.Author
		License     *loader.Link
		Theme       *loader.Theme

		RSS     *RSS
		Atom    *RSS
		Sitemap *Sitemap

		Uptime   time.Time
		Created  time.Time
		Modified time.Time
		Builded  time.Time // 最后次编译时间

		Tags     *Tags
		Index    *Index
		Archives *Archives
	}
)

// RSS Atom 和 RSS 的相关配置项
type RSS struct {
	*loader.RSS
	Permalink string
	XSL       string
}

// Sitemap 的相关配置项
type Sitemap struct {
	*loader.Sitemap
	Permalink string
	XSL       string
}

func newRSS(conf *loader.Config, rss *loader.RSS, path, xsl string) *RSS {
	r := &RSS{
		RSS:       rss,
		Permalink: buildURL(conf.URL, path),
	}

	if xsl != "" {
		r.XSL = buildThemeURL(conf.URL, conf.Theme, xsl)
	}

	return r
}

func newSitemap(conf *loader.Config, theme *loader.Theme) *Sitemap {
	sm := &Sitemap{
		Sitemap:   conf.Sitemap,
		Permalink: buildURL(conf.URL, vars.SitemapXML),
	}

	if theme.Sitemap != "" {
		sm.XSL = buildThemeURL(conf.URL, conf.Theme, theme.Sitemap)
	}

	return sm
}

// Load 加载并处理数据
func Load(dir string) (*Data, error) {
	conf, err := loader.LoadConfig(filepath.Join(dir, vars.ConfYAML))
	if err != nil {
		return nil, err
	}

	tags, err := loader.LoadTags(filepath.Join(dir, vars.TagsYAML))
	if err != nil {
		return nil, err
	}

	posts, err := loader.LoadPosts(dir)
	if err != nil {
		return nil, err
	}

	theme, err := loader.LoadTheme(dir, conf.Theme)
	if err != nil {
		return nil, err
	}

	return build(conf, tags, posts, theme)
}

func build(conf *loader.Config, tags *loader.Tags, posts []*loader.Post, theme *loader.Theme) (*Data, error) {
	var suffix string
	if conf.TitleSeparator != "" {
		suffix = conf.TitleSeparator + conf.Title
	}

	ts, err := buildTags(conf, tags)
	if err != nil {
		return nil, err
	}

	index, err := buildPosts(conf, theme, posts)
	if err != nil {
		return nil, err
	}

	archives, err := buildArchives(conf, index.Posts)
	if err != nil {
		return nil, err
	}

	created, modified, err := relationTagsPosts(ts.Tags, index.Posts)
	if err != nil {
		return nil, err
	}

	data := &Data{
		URL:         conf.URL,
		Title:       conf.Title,
		Subtitle:    conf.Subtitle,
		TitleSuffix: suffix,
		Icon:        conf.Icon,
		Language:    conf.Language,
		Authors:     conf.Authors,
		License:     conf.License,
		Theme:       theme,

		Uptime:   conf.Uptime,
		Builded:  time.Now(),
		Created:  created,
		Modified: modified,

		Tags:     ts,
		Index:    index,
		Archives: archives,
	}

	if conf.RSS != nil {
		data.RSS = newRSS(conf, conf.RSS, vars.RssXML, theme.RSS)
	}
	if conf.Atom != nil {
		data.Atom = newRSS(conf, conf.Atom, vars.AtomXML, theme.Atom)
	}
	if conf.Sitemap != nil {
		data.Sitemap = newSitemap(conf, theme)
	}

	return data, nil
}

func buildURL(url string, p ...string) string {
	if url == "" || url[len(url)-1] != '/' {
		url += "/"
	}

	pp := path.Join(p...)

	if len(pp) == 0 {
		return url
	}

	if pp[0] == '/' {
		return url + pp[1:]
	}
	return url + pp
}

func buildThemeURL(url, themeID string, p ...string) string {
	pp := make([]string, 0, len(p))
	pp = append(pp, vars.ThemesDir, themeID)
	return buildURL(url, append(pp, p...)...)
}

func buildPath(slug string) string {
	if slug == "" {
		panic("slug 不能为空")
	}

	if slug[0] == '/' || slug[0] == os.PathSeparator {
		slug = slug[1:]
	}
	return slug + ".xml"
}
