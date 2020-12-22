// SPDX-License-Identifier: MIT

package loader

import (
	"strconv"
	"time"

	"github.com/issue9/sliceutil"
)

// TimeFormat 所有时间的解析格式
const TimeFormat = time.RFC3339

// 归档的类型
const (
	ArchiveTypeYear  = "year"
	ArchiveTypeMonth = "month"
)

// 归档的排序方式
const (
	ArchiveOrderDesc = "desc"
	ArchiveOrderAsc  = "asc"
)

// Config 配置信息，用于从文件中读取
type Config struct {
	Title           string        `yaml:"title"`
	TitleSeparator  string        `yaml:"titleSeparator"`
	Language        string        `yaml:"language"`
	Subtitle        string        `yaml:"subtitle,omitempty"`
	Uptime          time.Time     `yaml:"uptime"`
	PageSize        int           `yaml:"pageSize"`
	Icon            *Icon         `yaml:"icon,omitempty"`
	Menus           []*Link       `yaml:"menus,omitempty"`
	Author          *Author       `yaml:"author"`
	License         *Link         `yaml:"license"`
	LongDateFormat  string        `yaml:"longDateFormat"`
	ShortDateFormat string        `yaml:"shortDateFormat"`
	Outdated        time.Duration `yaml:"outdated,omitempty"`
	Theme           string        `yaml:"theme"`

	Archive *Archive `yaml:"archive"`
	RSS     *RSS     `yaml:"rss,omitempty"`
	Atom    *RSS     `yaml:"atom,omitempty"`
	Sitemap *Sitemap `yaml:"sitemap,omitempty"`
}

// RSS RSS 和 Atom 相关的配置项
type RSS struct {
	Title string `yaml:"title"`
	Type  string `yaml:"type,omitempty"`
	Size  int    `yaml:"size"` // 显示数量
}

// Sitemap sitemap 相关的配置
type Sitemap struct {
	Priority   float64 `yaml:"priority"`            // 默认的优先级
	Changefreq string  `yaml:"changefreq"`          // 默认的更新频率
	EnableTag  bool    `yaml:"enableTag,omitempty"` // 是否将标签相关的页面写入 sitemap

	// 文章可以指定一个专门的值
	PostPriority   float64 `yaml:"postPriority"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

// Archive 存档页的配置内容
type Archive struct {
	Order  string `yaml:"order"`  // 排序方式
	Type   string `yaml:"type"`   // 存档的分类方式，可以按年或是按月
	Format string `yaml:"format"` // 标题的格式化字符串，被 time.Format 所格式化。
}

func loadConfig(path string) (*Config, error) {
	conf := &Config{}

	if err := loadYAML(path, conf); err != nil {
		return nil, err
	}
	if err := conf.sanitize(); err != nil {
		err.File = path
		return nil, err
	}

	return conf, nil
}

func (conf *Config) sanitize() *FieldError {
	if len(conf.Language) == 0 {
		conf.Language = "cmn-Hans"
	}

	if conf.PageSize <= 0 {
		return &FieldError{Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &FieldError{Message: "不能为空", Field: "longDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &FieldError{Message: "不能为空", Field: "shortDateFormat"}
	}

	if conf.Outdated < 0 {
		return &FieldError{Message: "必须大于 0", Field: "outdated"}
	}

	// icon
	if conf.Icon != nil {
		if err := conf.Icon.sanitize(); err != nil {
			err.Field = "icon." + err.Field
			return err
		}
	}

	// Author
	if conf.Author == nil {
		return &FieldError{Message: "不能为空", Field: "author"}
	}
	if err := conf.Author.sanitize(); err != nil {
		err.Field = "author." + err.Field
		return err
	}

	if len(conf.Title) == 0 {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	// theme
	if len(conf.Theme) == 0 {
		return &FieldError{Message: "不能为空", Field: "theme"}
	}

	// archive
	if conf.Archive == nil {
		return &FieldError{Message: "不能为空", Field: "archive"}
	}
	if err := conf.Archive.sanitize(); err != nil {
		err.Field = "archive." + err.Field
		return err
	}

	// license
	if conf.License == nil {
		return &FieldError{Message: "不能为空", Field: "license"}
	}
	if err := conf.License.sanitize(); err != nil {
		err.Field = "license." + err.Field
		return err
	}

	// rss
	if conf.RSS != nil {
		if err := conf.RSS.sanitize(conf); err != nil {
			err.Field = "rss." + err.Field
			return err
		}
	}

	// atom
	if conf.Atom != nil {
		if err := conf.Atom.sanitize(conf); err != nil {
			err.Field = "atom." + err.Field
			return err
		}
	}

	// sitemap
	if conf.Sitemap != nil {
		if err := conf.Sitemap.sanitize(); err != nil {
			err.Field = "sitemap." + err.Field
			return err
		}
	}

	// menus
	for index, link := range conf.Menus {
		if err := link.sanitize(); err != nil {
			err.Field = "menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

func (rss *RSS) sanitize(conf *Config) *FieldError {
	if rss.Size <= 0 {
		return &FieldError{Message: "必须大于 0", Field: "size"}
	}

	if len(rss.Title) == 0 {
		rss.Title = conf.Title
	}

	return nil
}

// 检测 sitemap 取值是否正确
func (s *Sitemap) sanitize() *FieldError {
	switch {
	case s.Priority > 1 || s.Priority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "priority"}
	case s.PostPriority > 1 || s.PostPriority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "postPriority"}
	case !inStrings(s.Changefreq, changereqs):
		return &FieldError{Message: "取值不正确", Field: "changefreq"}
	case !inStrings(s.PostChangefreq, changereqs):
		return &FieldError{Message: "取值不正确", Field: "postChangefreq"}
	}

	return nil
}

func (a *Archive) sanitize() *FieldError {
	if len(a.Type) == 0 {
		a.Type = ArchiveTypeYear
	} else {
		if a.Type != ArchiveTypeMonth && a.Type != ArchiveTypeYear {
			return &FieldError{Message: "取值不正确", Field: "type"}
		}
	}

	if len(a.Order) == 0 {
		a.Order = ArchiveOrderDesc
	} else {
		if a.Order != ArchiveOrderAsc && a.Order != ArchiveOrderDesc {
			return &FieldError{Message: "取值不正确", Field: "order"}
		}
	}

	if len(a.Format) == 0 {
		return &FieldError{Message: "不能为空", Field: "format"}
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
