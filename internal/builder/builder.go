// SPDX-License-Identifier: MIT

// Package builder 用于将由 laoder 加载的数据进行二次加工
package builder

import (
	"path/filepath"
	"time"

	"github.com/caixw/blogit/internal/loader"
)

// Build 打包目录下的内容
func Build(src string) (*Data, error) {
	conf, err := loader.LoadConfig(filepath.Join(src, "conf.yaml"))
	if err != nil {
		return nil, err
	}

	tags, err := loader.LoadTags(filepath.Join(src, "tags.yaml"))
	if err != nil {
		return nil, err
	}

	posts, err := loader.LoadPosts(src)
	if err != nil {
		return nil, err
	}

	return build(conf, tags, posts)
}

func build(conf *loader.Config, tags []*loader.Tag, posts []*loader.Post) (*Data, error) {
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

	ts, err := buildTags(conf, tags)
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

	data := &Data{
		Title:       conf.Title,
		Subtitle:    conf.Subtitle,
		TitleSuffix: suffix,
		Icon:        icon,
		Menus:       menus,
		Theme:       conf.Theme,

		Uptime:   conf.Uptime,
		Builded:  time.Now(),
		Created:  created,
		Modified: modified,

		Tags:  ts,
		Posts: ps,
	}

	return data, nil
}
