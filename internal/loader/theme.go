// SPDX-License-Identifier: MIT

package loader

import (
	"path/filepath"
	"strconv"

	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/internal/utils"
	"github.com/caixw/blogit/internal/vars"
)

// Theme 主题
type Theme struct {
	ID          string    `yaml:"-"`
	URL         string    `yaml:"url"`
	Description string    `yaml:"description,omitempty"`
	Authors     []*Author `yaml:"authors,omitempty"`
	Screenshots []string  `yaml:"screenshots,omitempty"`

	Templates []string `yaml:"templates,omitempty"`

	// 部分可选内容的模板，如果为空，则其输出相应的 xml 文件时不会为其添加 xsl 文件。
	// 模板名称为相对于当前主题目录的文件路径。
	Sitemap string `yaml:"sitemap,omitempty"`
	RSS     string `yaml:"rss,omitempty"`
	Atom    string `yaml:"atom,omitempty"`
}

// dir 为主题目录；id 为主题目录的名称
func (t *Theme) sanitize(dir, id string) *FieldError {
	t.ID = id

	for index, author := range t.Authors {
		if err := author.sanitize(); err != nil {
			err.Field = "authors[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	if sliceutil.Count(t.Templates, func(i int) bool { return t.Templates[i] == vars.DefaultTemplate }) == 0 {
		t.Templates = append(t.Templates, vars.DefaultTemplate)
	}
	indexes := sliceutil.Dup(t.Templates, func(i, j int) bool { return t.Templates[i] == t.Templates[j] })
	if len(indexes) > 0 {
		return &FieldError{Message: "重复的值模板列表", Field: "templates." + t.Templates[indexes[0]]}
	}

	for index, s := range t.Screenshots {
		if !utils.FileExists(filepath.Join(dir, s)) {
			return &FieldError{Message: "不存在的示例图", Field: "screenshots[" + strconv.Itoa(index) + "]"}
		}
	}
	indexes = sliceutil.Dup(t.Screenshots, func(i, j int) bool { return t.Screenshots[i] == t.Screenshots[j] })
	if len(indexes) > 0 {
		return &FieldError{Message: "重复的值示例图", Field: "screenshots[" + strconv.Itoa(indexes[0]) + "]"}
	}

	if t.Sitemap != "" {
		if !utils.FileExists(filepath.Join(dir, t.Sitemap)) {
			return &FieldError{Message: "不存在该模板文件", Field: "sitemap", Value: t.Sitemap}
		}
	}

	if t.RSS != "" {
		if !utils.FileExists(filepath.Join(dir, t.RSS)) {
			return &FieldError{Message: "不存在该模板文件", Field: "rss", Value: t.RSS}
		}
	}

	if t.Atom != "" {
		if !utils.FileExists(filepath.Join(dir, t.Atom)) {
			return &FieldError{Message: "不存在该模板文件", Field: "atom", Value: t.Atom}
		}
	}

	return nil
}

// LoadTheme 加载指定主题
func LoadTheme(dir, name string) (*Theme, error) {
	dir = filepath.Join(dir, vars.ThemesDir, name)
	path := filepath.Join(dir, vars.ThemeYAML)

	theme := &Theme{}
	if err := loadYAML(path, &theme); err != nil {
		return nil, err
	}

	if err := theme.sanitize(dir, name); err != nil {
		err.File = path
		return nil, err
	}
	return theme, nil
}
