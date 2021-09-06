// SPDX-License-Identifier: MIT

// Package data 对加载的数据进行二次加工
package data

import (
	"io/fs"
	"path"
	"time"

	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
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
		Theme       *Theme
		Highlights  []*Highlight
		Menus       []*loader.Link

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
		Posts    []*Post
		Indexes  []*Index
		Archives *Archives
	}
)

// Load 加载并处理数据
//
// preview 表示是否为预览模式，在预览模式下会加载草稿同；
// 如果 baseURL 不为空，则会替换配置文件中的 URL 字段。
func Load(fs fs.FS, preview bool, baseURL string) (*Data, error) {
	conf, err := loader.LoadConfig(fs, vars.ConfYAML)
	if err != nil {
		return nil, err
	}
	if baseURL != "" {
		conf.URL = baseURL
	}

	tags, err := loader.LoadTags(fs, vars.TagsYAML)
	if err != nil {
		return nil, err
	}

	posts, err := loader.LoadPosts(fs, preview)
	if err != nil {
		return nil, err
	}

	theme, err := loader.LoadTheme(fs, conf.Theme)
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

	ps, err := buildPosts(conf, theme, posts)
	if err != nil {
		return nil, err
	}

	archives, err := buildArchives(conf, ps)
	if err != nil {
		return nil, err
	}

	created, modified, err := relationTagsPosts(ts.Tags, ps)
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
		Theme:       newTheme(theme),
		Highlights:  newHighlights(conf, theme),
		Menus:       conf.Menus,

		Uptime:   conf.Uptime,
		Builded:  time.Now(),
		Created:  created,
		Modified: modified,

		Tags:     ts,
		Posts:    ps,
		Indexes:  buildIndexes(conf, ps),
		Archives: archives,
	}

	// 获得一份按时间排序的列表，诸如 rss 等不应该受自定义排序的影响，始终以时间作为排序。
	sorted := sortPostsByCreated(ps)

	if conf.RSS != nil {
		data.RSS = newRSS(conf, conf.RSS, vars.RssXML, theme.RSS, sorted)
	}
	if conf.Atom != nil {
		data.Atom = newRSS(conf, conf.Atom, vars.AtomXML, theme.Atom, sorted)
	}
	if conf.Sitemap != nil {
		data.Sitemap = newSitemap(conf, theme)
	}
	if conf.Robots != nil {
		data.Robots = newRobots(conf, data.Sitemap)
	}
	if conf.Profile != nil {
		data.Profile = newProfile(conf, sorted)
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
