// SPDX-License-Identifier: MIT

package loader

import (
	"strings"
	"unicode"

	"github.com/issue9/localeutil"
)

// Profile 用于生成 github.com/profile 中的 README.md 内容
type Profile struct {
	// 以下字段提供了类似于以下格式的 markdown 内容：
	//  ### Title
	//
	//  post1
	//  post2
	//  post3
	//
	//  ##### Footer
	//
	// Title 和 Footer 的前缀 # 是固定的，不需要用户给字，即使用户给了，也会被删除。
	Title  string `yaml:"title"`
	Footer string `yaml:"footer"` // 页脚
	Size   int    `yaml:"size"`   // 显示最近添加的文章条数
}

func (p *Profile) sanitize() *FieldError {
	if p.Title == "" {
		return &FieldError{Field: "title", Message: localeutil.Phrase("can not be empty")}
	}
	p.Title = "### " + trimHeadPrefix(p.Title)

	if p.Size <= 0 {
		return &FieldError{Field: "size", Message: localeutil.Phrase("should great than zero")}
	}

	if p.Footer != "" {
		p.Footer = "##### " + trimHeadPrefix(p.Footer)
	}

	return nil
}

func trimHeadPrefix(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || r == '#'
	})
}
