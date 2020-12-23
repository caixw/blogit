// SPDX-License-Identifier: MIT

package loader

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/issue9/sliceutil"
)

// Tag 描述标签信息
type Tag struct {
	Title     string `yaml:"title"`
	HTMLTitle string `yaml:"-"`

	Slug    string `yaml:"slug"`            // 唯一名称
	Color   string `yaml:"color,omitempty"` // 标签颜色。若未指定，则继承父容器
	Content string `yaml:"content"`         // 对该标签的详细描述

	Permalink string    `yaml:"-"`
	Posts     int       `yaml:"-"`
	Created   time.Time `yaml:"-"`
	Modified  time.Time `yaml:"-"`
}

func (data *Data) loadTags(filename string) error {
	tags := make([]*Tag, 0, 100)
	path := filepath.Join(data.Dir, filename)

	if err := loadYAML(path, &tags); err != nil {
		return err
	}

	for index, tag := range tags {
		if err := tag.sanitize(tags, data.Config); err != nil {
			err.File = path
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	data.Tags = tags
	return nil
}

func (data *Data) checkTags() error {
	for _, p := range data.Posts {
		p.Tags = make([]*Tag, 0, len(p.TagString))

		for _, tag := range p.TagString {
			t := data.findTagByName(tag)
			if t == nil {
				return &FieldError{File: p.Slug, Message: "不存在", Field: "tags." + tag}
			}
			p.Tags = append(p.Tags, t)

			if t.Created.Before(p.Created) {
				t.Created = p.Created
			}

			if t.Modified.Before(p.Modified) {
				t.Modified = p.Modified
			}
		}
	}
	return nil
}

func (data *Data) findTagByName(slug string) *Tag {
	for _, t := range data.Tags {
		if t.Slug == slug {
			return t
		}
	}
	return nil
}

func (tag *Tag) sanitize(tags []*Tag, conf *Config) *FieldError {
	if len(tag.Slug) == 0 {
		return &FieldError{Message: "不能为空", Field: "slug"}
	}

	if len(tag.Title) == 0 {
		return &FieldError{Message: "不能为空", Field: "title"}
	}
	tag.HTMLTitle = tag.Title + conf.titleSuffix

	if len(tag.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "content"}
	}

	cnt := sliceutil.Count(tags, func(i int) bool {
		return tags[i].Slug == tag.Slug
	})
	if cnt > 1 {
		return &FieldError{Message: "重复的值", Field: "slug"}
	}

	tag.Permalink = conf.BuildURL(tag.Slug + ".xml")

	return nil
}
