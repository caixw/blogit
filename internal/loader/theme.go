// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"path"
	"strconv"

	"github.com/alecthomas/chroma/styles"
	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/internal/filesystem"
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

	// 指定高亮的主题名称
	//
	// 名称值可以从 https://pkg.go.dev/github.com/alecthomas/chroma@v0.8.2/styles 获取
	Highlights []*Highlight `yaml:"highlights,omitempty"`

	// 部分可选内容的模板，如果为空，则其输出相应的 xml 文件时不会为其添加 xsl 文件。
	// 模板名称为相对于当前主题目录的文件路径。
	Sitemap string `yaml:"sitemap,omitempty"`
	RSS     string `yaml:"rss,omitempty"`
	Atom    string `yaml:"atom,omitempty"`
}

// Highlight 高亮主题指定
type Highlight struct {
	Name  string `yaml:"name"` // 指向的主题名称
	Media string `yaml:"media,omitempty"`
}

// dir 为当前主题所在的目录；
// id 为主题目录的名称
func (t *Theme) sanitize(fs fs.FS, dir, id string) *FieldError {
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
		if !filesystem.Exists(fs, path.Join(dir, s)) {
			return &FieldError{Message: "不存在的示例图", Field: "screenshots[" + strconv.Itoa(index) + "]"}
		}
	}
	indexes = sliceutil.Dup(t.Screenshots, func(i, j int) bool { return t.Screenshots[i] == t.Screenshots[j] })
	if len(indexes) > 0 {
		return &FieldError{Message: "重复的示例图", Field: "screenshots[" + strconv.Itoa(indexes[0]) + "]"}
	}

	if t.Sitemap != "" && !filesystem.Exists(fs, path.Join(dir, t.Sitemap)) {
		return &FieldError{Message: "不存在该模板文件", Field: "sitemap", Value: t.Sitemap}
	}

	if t.RSS != "" && !filesystem.Exists(fs, path.Join(dir, t.RSS)) {
		return &FieldError{Message: "不存在该模板文件", Field: "rss", Value: t.RSS}
	}

	if t.Atom != "" && !filesystem.Exists(fs, path.Join(dir, t.Atom)) {
		return &FieldError{Message: "不存在该模板文件", Field: "atom", Value: t.Atom}
	}

	var mediaIsEmpty bool
	for index, h := range t.Highlights {
		i := strconv.Itoa(index)
		prefix := "highlight[" + i + "]."

		if err := h.sanitize(); err != nil {
			err.Field = prefix + err.Field
			return err
		}

		if h.Media == "" {
			if mediaIsEmpty {
				return &FieldError{Message: "只能一个为空", Field: prefix + "media"}
			}
			mediaIsEmpty = true
		}
	}

	return nil
}

var highlightCSSName = styles.Names()

func (h *Highlight) sanitize() *FieldError {
	names := highlightCSSName
	if sliceutil.Count(names, func(i int) bool { return names[i] == h.Name }) == 0 {
		return &FieldError{Message: "不存在", Field: "name", Value: h.Name}
	}

	return nil
}

// LoadTheme 加载指定主题
func LoadTheme(fs fs.FS, id string) (*Theme, error) {
	dir := path.Join(vars.ThemesDir, id)
	p := path.Join(dir, vars.ThemeYAML)

	theme := &Theme{}
	if err := loadYAML(fs, p, &theme); err != nil {
		return nil, err
	}

	if err := theme.sanitize(fs, dir, id); err != nil {
		err.File = p
		return nil, err
	}
	return theme, nil
}
