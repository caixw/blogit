// SPDX-License-Identifier: MIT

// Package data 对加载的数据进行二次加工
package data

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/caixw/blogit/internal/loader"
)

type (
	// Data 处理后的数据
	Data struct {
		URL         string
		Title       string
		Subtitle    string
		TitleSuffix string // 每篇文章标题的后缀
		Icon        *loader.Icon
		Menus       []*loader.Menu
		Language    string
		Authors     []*loader.Author
		License     *loader.License
		Theme       *loader.Theme

		Archive *loader.Archive
		RSS     *loader.RSS
		Atom    *loader.RSS
		Sitemap *loader.Sitemap

		LongDateFormat  string
		ShortDateFormat string
		Uptime          time.Time
		Created         time.Time
		Modified        time.Time
		Builded         time.Time // 最后次编译时间

		Tags  []*Tag
		Posts []*Post
	}
)

// Load 加载并处理数据
func Load(dir string) (*Data, error) {
	conf, err := loader.LoadConfig(filepath.Join(dir, "conf.yaml"))
	if err != nil {
		return nil, err
	}

	tags, err := loader.LoadTags(filepath.Join(dir, "tags.yaml"))
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

func build(conf *loader.Config, tags []*loader.Tag, posts []*loader.Post, theme *loader.Theme) (*Data, error) {
	var suffix string
	if conf.TitleSeparator != "" {
		suffix = conf.TitleSeparator + conf.Title
	}

	ts, err := buildTags(tags)
	if err != nil {
		return nil, err
	}

	ps, err := buildPosts(conf, theme, posts)
	if err != nil {
		return nil, err
	}

	created, modified, err := checkTags(ts, ps)
	if err != nil {
		return nil, err
	}

	data := &Data{
		URL:         conf.URL,
		Title:       conf.Title,
		Subtitle:    conf.Subtitle,
		TitleSuffix: suffix,
		Icon:        conf.Icon,
		Menus:       conf.Menus,
		Language:    conf.Language,
		Authors:     conf.Authors,
		License:     conf.License,
		Theme:       theme,

		Archive: conf.Archive,
		RSS:     conf.RSS,
		Atom:    conf.Atom,
		Sitemap: conf.Sitemap,

		LongDateFormat:  conf.LongDateFormat,
		ShortDateFormat: conf.ShortDateFormat,
		Uptime:          conf.Uptime,
		Builded:         time.Now(),
		Created:         created,
		Modified:        modified,

		Tags:  ts,
		Posts: ps,
	}

	return data, nil
}

// BuildURL 根据配置网站域名生成地址
func (data *Data) BuildURL(p ...string) string {
	pp := path.Join(p...)

	if len(pp) == 0 {
		return data.URL
	}

	if pp[0] == '/' {
		return data.URL + pp[1:]
	}
	return data.URL + pp
}

// BuildThemeURL 根据配置网站域名生成主题下的文件地址
func (data *Data) BuildThemeURL(p ...string) string {
	pp := make([]string, 0, len(p)+2)
	pp = append(pp, "themes", data.Theme.ID)
	return data.BuildURL(append(pp, p...)...)
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
