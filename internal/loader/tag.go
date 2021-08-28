// SPDX-License-Identifier: MIT

package loader

import (
	"bytes"
	"io/fs"
	"strconv"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"
)

const (
	TagOrderTypeSize    = "size" // 按关联的文章数量排序
	TagOrderTypeDefault = ""     // 默认方式，即文件内容的顺序
)

// Tags 标签页的数据
type Tags struct {
	Title       string `yaml:"title"` // 存档页的标题
	Keywords    string `yaml:"keywords,omitempty"`
	Description string `yaml:"description,omitempty"`
	Order       string `yaml:"order,omitempty"` // 排序方式
	OrderType   string `yaml:"orderType,omitempty"`
	Tags        []*Tag `yaml:"tags,omitempty"`
}

// Tag 描述标签信息
type Tag struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"` // 对该标签的详细描述
	Slug    string `yaml:"slug"`    // 唯一名称
}

// LoadTags 加载标签列表
func LoadTags(fs fs.FS, path string) (*Tags, error) {
	tags := &Tags{}
	if err := loadYAML(fs, path, &tags); err != nil {
		return nil, err
	}

	if err := tags.sanitize(); err != nil {
		err.File = path
		return nil, err
	}

	return tags, nil
}

func (tags *Tags) sanitize() *FieldError {
	if tags.Title == "" {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "title"}
	}

	switch tags.Order {
	case "":
		tags.Order = OrderDesc
	case OrderAsc, OrderDesc:
	default:
		return &FieldError{Message: localeutil.Phrase("invalid value"), Field: "order"}
	}

	switch tags.OrderType {
	case TagOrderTypeDefault, TagOrderTypeSize:
	default:
		return &FieldError{Message: localeutil.Phrase("invalid value"), Field: "orderType"}
	}

	for index, tag := range tags.Tags {
		if err := tag.sanitize(tags); err != nil {
			err.Field = "tag[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

func (tag *Tag) sanitize(tags *Tags) *FieldError {
	if len(tag.Slug) == 0 {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "slug"}
	}

	if len(tag.Title) == 0 {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "title"}
	}

	if len(tag.Content) == 0 {
		return &FieldError{Message: localeutil.Phrase("can not be empty"), Field: "content"}
	}

	// 将 markdown 转换成 html
	buf := new(bytes.Buffer)
	if err := markdown.Convert([]byte(tag.Content), buf); err != nil {
		return &FieldError{Message: localeutil.Phrase(err.Error()), Field: "content", Value: tag.Content}
	}
	tag.Content = buf.String()

	if tags != nil && tags.Tags != nil {
		cnt := sliceutil.Count(tags.Tags, func(i int) bool {
			return tags.Tags[i].Slug == tag.Slug
		})
		if cnt > 1 {
			return &FieldError{Message: localeutil.Phrase("duplicate value"), Field: "slug", Value: tag.Slug}
		}
	}

	return nil
}
