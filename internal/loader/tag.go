// SPDX-License-Identifier: MIT

package loader

import (
	"bytes"
	"strconv"

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
func LoadTags(path string) (*Tags, error) {
	tags := &Tags{}
	if err := loadYAML(path, &tags); err != nil {
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
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	switch tags.Order {
	case "":
		tags.Order = OrderDesc
	case OrderAsc, OrderDesc:
	default:
		return &FieldError{Message: "无效的值", Field: "order"}
	}

	switch tags.OrderType {
	case TagOrderTypeDefault, TagOrderTypeSize:
	default:
		return &FieldError{Message: "无效的值", Field: "orderType"}
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
		return &FieldError{Message: "不能为空", Field: "slug"}
	}

	if len(tag.Title) == 0 {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	if len(tag.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "content"}
	}

	// 将 markdown 转换成 html
	buf := new(bytes.Buffer)
	if err := markdown.Convert([]byte(tag.Content), buf); err != nil {
		return &FieldError{Message: err.Error(), Field: "content", Value: tag.Content}
	}
	tag.Content = buf.String()

	if tags != nil && tags.Tags != nil {
		cnt := sliceutil.Count(tags.Tags, func(i int) bool {
			return tags.Tags[i].Slug == tag.Slug
		})
		if cnt > 1 {
			return &FieldError{Message: "重复的值", Field: "slug", Value: tag.Slug}
		}
	}

	return nil
}
