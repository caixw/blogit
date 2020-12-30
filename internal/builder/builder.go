// SPDX-License-Identifier: MIT

// Package builder 用于将由 laoder 加载的数据进行二次加工
package builder

import (
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
		Icon        *Icon
		Menus       []*Menu
		Theme       *Theme

		LongDateFormat  string
		ShortDateFormat string
		Uptime          time.Time
		Created         time.Time
		Modified        time.Time
		Builded         time.Time // 最后次编译时间

		Tags     []*Tag
		Posts    []*Post
		Archives []*Archive
	}

	// Icon 图标信息
	Icon = loader.Icon

	// License 表示链接信息
	License = loader.License

	// Menu 采单项
	Menu = loader.Menu

	// Author 表示作者信息
	Author = loader.Author

	// Theme 主题
	Theme = loader.Theme
)

// Build 打包目录下的内容
func Build(dir string) (*Data, error) {
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

	icon := &Icon{
		URL:   conf.Icon.URL,
		Type:  conf.Icon.Type,
		Sizes: conf.Icon.Sizes,
	}

	menus := make([]*Menu, 0, len(conf.Menus))
	for _, m := range conf.Menus {
		menus = append(menus, &Menu{
			Icon:  m.Icon,
			Title: m.Title,
			URL:   m.URL,
			Text:  m.Text,
		})
	}

	ts, err := buildTags(tags)
	if err != nil {
		return nil, err
	}

	ps, err := buildPosts(conf, posts)
	if err != nil {
		return nil, err
	}

	created, modified, err := checkTags(ts, ps)
	if err != nil {
		return nil, err
	}

	archives, err := buildArchives(conf, ps)
	if err != nil {
		return nil, err
	}

	data := &Data{
		URL: conf.URL,

		Title:       conf.Title,
		Subtitle:    conf.Subtitle,
		TitleSuffix: suffix,
		Icon:        icon,
		Menus:       menus,
		Theme:       theme,

		LongDateFormat:  conf.LongDateFormat,
		ShortDateFormat: conf.ShortDateFormat,
		Uptime:          conf.Uptime,
		Builded:         time.Now(),
		Created:         created,
		Modified:        modified,

		Tags:     ts,
		Posts:    ps,
		Archives: archives,
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

	return data.BuildThemeURL(append(pp, p...)...)
}
