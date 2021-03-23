// SPDX-License-Identifier: MIT

// Package loader 加载数据内容
//
// 仅加载各个模块的自身的数据，并判断格式是否正确。
// 但是不会对各个模块之间的关联数据进行校验。
package loader

import (
	"fmt"
	"io/ioutil"
	"mime"
	"path"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/issue9/validation/is"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"gopkg.in/yaml.v2"
)

var markdown = goldmark.New(goldmark.WithExtensions(
	extension.GFM,
	extension.Strikethrough,
	extension.Footnote,
	meta.Meta,
	highlighting.NewHighlighting(
		highlighting.WithStyle("monokai"),
		highlighting.WithFormatOptions(
			html.WithLineNumbers(true),
		),
	),
))

// 排序方式
const (
	OrderDesc = "desc"
	OrderAsc  = "asc"
)

// FieldError 表示配置项内容的错误信息
type FieldError struct {
	File    string
	Field   string
	Message string
	Value   interface{}
}

// Link 描述链接的内容
type Link struct {
	URL  string `yaml:"url"`  // 链接地址
	Text string `yaml:"text"` // 链接的文本
}

// Icon 表示网站图标，比如 html>head>link.rel="short icon"
type Icon struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type"` // mime type
	Sizes string `yaml:"sizes"`
}

// Author 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

func (err *FieldError) Error() string {
	if err.Value == nil {
		return fmt.Sprintf("%s 位于 %s:%s", err.Message, err.File, err.Field)
	}
	return fmt.Sprintf("%s 位于 %s:%s，实际值为:%#v", err.Message, err.File, err.Field, err.Value)
}

func loadYAML(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (icon *Icon) sanitize() *FieldError {
	if len(icon.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	if icon.Type == "" {
		icon.Type = mime.TypeByExtension(path.Ext(icon.URL))
	}

	return nil
}

func (l *Link) sanitize() *FieldError {
	if len(l.Text) == 0 {
		return &FieldError{Field: "text", Message: "不能为空"}
	}

	if len(l.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (author *Author) sanitize() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "name", Message: "不能为空"}
	}

	if len(author.URL) > 0 && !is.URL(author.URL) {
		return &FieldError{Field: "url", Message: "不是一个正确的 URL", Value: author.URL}
	}

	if len(author.Avatar) > 0 && !is.URL(author.Avatar) {
		return &FieldError{Field: "avatar", Message: "不是一个正确的 URL", Value: author.Avatar}
	}

	if len(author.Email) > 0 && !is.Email(author.Email) {
		return &FieldError{Field: "email", Message: "不是一个正确的 Email", Value: author.Email}
	}

	return nil
}
