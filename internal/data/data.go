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
		Author      *loader.Author
		License     *loader.Link
		Theme       *loader.Theme

		RSS     *RSS
		Atom    *RSS
		Sitemap *Sitemap
		Robots  *Robots
		Profile *Profile

		Uptime   time.Time
		Created  time.Time
		Modified time.Time
		Builded  time.Time // 最后次编译时间

		Tags     *Tags
		Index    *Index
		Archives *Archives
	}
)

// Load 加载并处理数据
func Load(dir, baseURL string) (*Data, error) {
	conf, err := loader.LoadConfig(filepath.Join(dir, vars.ConfYAML))
	if err != nil {
		return nil, err
	}
	if baseURL != "" {
		conf.URL = baseURL
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
		Author:      conf.Author,
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
		data.RSS = newRSS(conf, conf.RSS, vars.RssXML, theme.RSS, index.Posts)
	}
	if conf.Atom != nil {
		data.Atom = newRSS(conf, conf.Atom, vars.AtomXML, theme.Atom, index.Posts)
	}
	if conf.Sitemap != nil {
		data.Sitemap = newSitemap(conf, theme)
	}
	if conf.Robots != nil {
		data.Robots = newRobots(conf, data.Sitemap)
	}
	if conf.Profile != nil {
		data.Profile = newProfile(conf, index.Posts)
	}

	return data, nil
}

// BuildURL 将 p 添加到 baseURL 形成一条完整的 URL
func BuildURL(baseURL string, p ...string) string {
	if baseURL == "" || baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	pp := path.Join(p...)

	if len(pp) == 0 {
		return baseURL
	}

	if pp[0] == '/' {
		return baseURL + pp[1:]
	}
	return baseURL + pp
}

func buildThemeURL(baseURL, themeID string, p ...string) string {
	pp := make([]string, 0, len(p))
	pp = append(pp, vars.ThemesDir, themeID)
	return BuildURL(baseURL, append(pp, p...)...)
}

// 如果 slug 不再扩展名，再会加上默认的扩展名 .html
func buildPath(slug string) string {
	if slug == "" {
		panic("slug 不能为空")
	}

	if slug[0] == '/' || slug[0] == os.PathSeparator {
		slug = slug[1:]
	}

	if filepath.Ext(slug) != "" {
		return slug
	}
	return slug + vars.Ext
}
