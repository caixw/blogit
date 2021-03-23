// SPDX-License-Identifier: MIT

package loader

import "github.com/issue9/sliceutil"

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
	switch {
	case s.Title == "":
		return &FieldError{Message: "不能为空", Field: "title"}
	case s.Priority > 1 || s.Priority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "priority", Value: s.Priority}
	case s.PostPriority > 1 || s.PostPriority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "postPriority", Value: s.PostPriority}
	case !inStrings(s.Changefreq, changereqs):
		return &FieldError{Message: "取值不正确", Field: "changefreq", Value: s.Changefreq}
	case !inStrings(s.PostChangefreq, changereqs):
		return &FieldError{Message: "取值不正确", Field: "postChangefreq", Value: s.PostChangefreq}
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

func inStrings(val string, vals []string) bool {
	return sliceutil.Count(vals, func(i int) bool {
		return vals[i] == val
	}) > 0
}
