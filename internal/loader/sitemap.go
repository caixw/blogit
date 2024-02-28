// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package loader

import (
	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"
)

// Sitemap sitemap 相关的配置
type Sitemap struct {
	Title string `yaml:"title"`

	Priority   float64 `yaml:"priority"`            // 默认的优先级
	Changefreq string  `yaml:"changefreq"`          // 默认的更新频率
	EnableTag  bool    `yaml:"enableTag,omitempty"` // 是否将标签相关的页面写入 sitemap

	// 文章可以指定一个专门的值
	PostPriority   float64 `yaml:"postPriority"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

func (s *Sitemap) sanitize() *FieldError {
	chkPriority := func(v float64, field string) *FieldError {
		if v > 1 || v < 0 {
			return &FieldError{Message: localeutil.StringPhrase("should be float"), Field: field, Value: v}
		}
		return nil
	}
	chkChangefreq := func(v, field string) *FieldError {
		if !inStrings(v, changereqs) {
			return &FieldError{Message: InvalidValue, Field: field, Value: v}
		}
		return nil
	}

	if s.Title == "" {
		return &FieldError{Message: Required, Field: "title"}
	}

	if err := chkPriority(s.Priority, "priority"); err != nil {
		return err
	}
	if err := chkPriority(s.PostPriority, "postPriority"); err != nil {
		return err
	}
	if err := chkChangefreq(s.Changefreq, "changefreq"); err != nil {
		return err
	}
	if err := chkChangefreq(s.PostChangefreq, "postChangefreq"); err != nil {
		return err
	}

	return nil
}

var changereqs = []string{
	"never",
	"yearly",
	"monthly",
	"weekly",
	"daily",
	"hourly",
	"always",
}

func inStrings(val string, values []string) bool {
	return sliceutil.Count(values, func(s string, _ int) bool { return s == val }) > 0
}
