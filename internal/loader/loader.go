// SPDX-License-Identifier: MIT

// Package loader 加载数据内容
//
// 仅加载各个模块的自身的数据，并判断格式是否正确。
// 但是不会对各个模块之间的关联数据进行校验。
package loader

import (
	"io/fs"
	"mime"
	"path"

	"github.com/issue9/localeutil"
	"github.com/issue9/validation/is"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

const (
	Required     = localeutil.StringPhrase("can not be empty")
	GreatZero    = localeutil.StringPhrase("should great than zero")
	InvalidValue = localeutil.StringPhrase("invalid value")
	InvalidURL   = localeutil.StringPhrase("invalid url")
	NotFound     = localeutil.StringPhrase("not found")
	DupValue     = localeutil.StringPhrase("duplicate value")
)

// 排序方式
const (
	OrderDesc = "desc"
	OrderAsc  = "asc"
)

// FieldError 表示配置项内容的错误信息
type FieldError struct {
	File    string
	Field   string
	Message localeutil.Stringer
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
	return err.Message.LocaleString(message.NewPrinter(language.Und))
}

func (err *FieldError) LocaleString(p *message.Printer) string {
	if err.Value == nil {
		return localeutil.Phrase("%s at %s:%d", err.Message.LocaleString(p), err.File, err.Field).LocaleString(p)
	}
	return localeutil.Phrase("%s at %s:%d,value is %s", err.Message.LocaleString(p), err.File, err.Field, err.Value).LocaleString(p)
}

func loadYAML(f fs.FS, path string, v interface{}) error {
	data, err := fs.ReadFile(f, path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

func (icon *Icon) sanitize() *FieldError {
	if len(icon.URL) == 0 {
		return &FieldError{Field: "url", Message: Required}
	}

	if icon.Type == "" {
		icon.Type = mime.TypeByExtension(path.Ext(icon.URL))
	}

	return nil
}

func (l *Link) sanitize() *FieldError {
	if len(l.Text) == 0 {
		return &FieldError{Field: "text", Message: Required}
	}

	if len(l.URL) == 0 {
		return &FieldError{Field: "url", Message: Required}
	}

	return nil
}

func (author *Author) sanitize() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "name", Message: Required}
	}

	if len(author.URL) > 0 && !is.URL(author.URL) {
		return &FieldError{Field: "url", Message: InvalidURL, Value: author.URL}
	}

	if len(author.Avatar) > 0 && !is.URL(author.Avatar) {
		return &FieldError{Field: "avatar", Message: InvalidURL, Value: author.Avatar}
	}

	if len(author.Email) > 0 && !is.Email(author.Email) {
		return &FieldError{Field: "email", Message: InvalidURL, Value: author.Email}
	}

	return nil
}
