// SPDX-License-Identifier: MIT

package loader

import (
	"strings"
	"unicode"
)

// Profile 用于生成 github.com/profile 中的 README.md 内容
type Profile struct {
	Alternate string `yaml:"alternate"` // 采用此文件的内容代替

	// 当 alternate 为空，以下值才生效
	//
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
	if p.Alternate != "" {
		switch {
		case p.Title != "":
			return &FieldError{Field: "title", Message: "只能为空"}
		case p.Size != 0:
			return &FieldError{Field: "size", Message: "只能为空"}
		case p.Footer != "":
			return &FieldError{Field: "footer", Message: "只能为空"}
		}
		return nil
	}

	if p.Title == "" {
		return &FieldError{Field: "title", Message: "不能为空"}
	}
	p.Title = "### " + trimHeadPrefix(p.Title)

	if p.Size <= 0 {
		return &FieldError{Field: "size", Message: "必须大于 0"}
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
