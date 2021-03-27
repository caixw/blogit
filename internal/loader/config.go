// SPDX-License-Identifier: MIT

package loader

import (
	"strconv"
	"time"

	"github.com/issue9/validation/is"
)

// Config 配置信息，用于从文件中读取
type Config struct {
	Title    string `yaml:"title"`
	Subtitle string `yaml:"subtitle,omitempty"`

	// 标题后缀分隔符，文章页面浏览器标题上会加上此后缀，如果为空，则表示不需要后缀。
	TitleSeparator string `yaml:"titleSeparator,omitempty"`

	URL         string    `yaml:"url"` // 网站根域名，比如 https://example.com/blog
	Language    string    `yaml:"language,omitempty"`
	Uptime      time.Time `yaml:"uptime"`
	Icon        *Icon     `yaml:"icon,omitempty"`
	Author      *Author   `yaml:"author"` // 网站作者，在文章没有指定作者时，也采用此值。
	License     *Link     `yaml:"license"`
	Theme       string    `yaml:"theme"`
	Keywords    string    `yaml:"keywords,omitempty"`    // 所有页面默认情况下的 keywords
	Description string    `yaml:"description,omitempty"` // 所有页面默认情况下的 description

	Archive *Archive `yaml:"archive,omitempty"`
	RSS     *RSS     `yaml:"rss,omitempty"`
	Atom    *RSS     `yaml:"atom,omitempty"`
	Sitemap *Sitemap `yaml:"sitemap,omitempty"`
	Robots  []*Agent `yaml:"robots,omitempty"`  // 不为空，表示托管 robots.txt 的生成
	Profile *Profile `yaml:"profile,omitempty"` // 不为空，表示托管 README.md 的生成
}

// RSS RSS 和 Atom 相关的配置项
type RSS struct {
	Title string `yaml:"title,omitempty"`
	Size  int    `yaml:"size"` // 显示数量
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
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
	if len(conf.URL) == 0 || !is.URL(conf.URL) {
		return &FieldError{Message: "格式不正确", Field: "url", Value: conf.URL}
	}

	if len(conf.Language) == 0 {
		conf.Language = "cmn-Hans"
	}

	if conf.Uptime.IsZero() {
		return &FieldError{Message: "不能为空", Field: "uptime"}
	}

	// icon
	if conf.Icon != nil {
		if err := conf.Icon.sanitize(); err != nil {
			err.Field = "icon." + err.Field
			return err
		}
	}

	// Authors
	if conf.Author == nil {
		return &FieldError{Message: "不能为空", Field: "authors"}
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
		if err := conf.RSS.sanitize(); err != nil {
			err.Field = "rss." + err.Field
			return err
		}
	}

	// atom
	if conf.Atom != nil {
		if err := conf.Atom.sanitize(); err != nil {
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

	// robots.txt
	if conf.Robots != nil {
		for index, agent := range conf.Robots {
			if err := agent.sanitize(); err != nil {
				err.Field = "robots.[" + strconv.Itoa(index) + "]." + err.Field
				return err
			}
		}
	}

	// profile
	if conf.Profile != nil {
		if err := conf.Profile.sanitize(); err != nil {
			err.Field = "profile." + err.Field
			return err
		}
	}

	return nil
}

func (rss *RSS) sanitize() *FieldError {
	if rss.Title == "" {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	if rss.Size <= 0 {
		return &FieldError{Message: "必须大于 0", Field: "size", Value: rss.Size}
	}

	return nil
}
