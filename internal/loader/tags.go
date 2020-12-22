// SPDX-License-Identifier: MIT

package loader

import (
	"strconv"

	"github.com/issue9/sliceutil"
)

// Tag 描述标签信息
type Tag struct {
	Slug    string `yaml:"slug"`            // 唯一名称
	Title   string `yaml:"title"`           // 名称
	Color   string `yaml:"color,omitempty"` // 标签颜色。若未指定，则继承父容器
	Content string `yaml:"content"`         // 对该标签的详细描述
}

func loadTags(path string) ([]*Tag, error) {
	tags := make([]*Tag, 0, 100)
	if err := loadYAML(path, &tags); err != nil {
		return nil, err
	}

	for index, tag := range tags {
		if err := tag.sanitize(tags); err != nil {
			err.File = path
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return nil, err
		}
	}

	return tags, nil
}

func (tag *Tag) sanitize(tags []*Tag) *FieldError {
	if len(tag.Slug) == 0 {
		return &FieldError{Message: "不能为空", Field: "slug"}
	}

	if len(tag.Title) == 0 {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	if len(tag.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "content"}
	}

	cnt := sliceutil.Count(tags, func(i int) bool {
		return tags[i].Slug == tag.Slug
	})
	if cnt > 1 {
		return &FieldError{Message: "重复的值", Field: "slug"}
	}

	return nil
}
