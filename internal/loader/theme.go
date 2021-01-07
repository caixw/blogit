// SPDX-License-Identifier: MIT

package loader

import (
	"path/filepath"
	"strconv"

	"github.com/caixw/blogit/internal/utils"
	"github.com/issue9/sliceutil"
)

// Theme 主题
type Theme struct {
	ID          string    `yaml:"-"`
	Description string    `yaml:"description,omitempty"`
	Authors     []*Author `yaml:"authors,omitempty"`
	Templates   []string  `yaml:"templates,omitempty"`
	Screenshots []string  `yaml:"screenshots,omitempty"`
}

func (t *Theme) sanitize(dir, id string) *FieldError {
	t.ID = id

	for index, author := range t.Authors {
		if err := author.sanitize(); err != nil {
			err.Field = "authors[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	for _, tpl := range t.Templates {
		if !utils.FileExists(filepath.Join(dir, tpl)) {
			return &FieldError{Message: "不存在该模板文件", Field: "templates." + tpl}
		}
	}
	indexes := sliceutil.Dup(t.Templates, func(i, j int) bool { return t.Templates[i] == t.Templates[j] })
	if len(indexes) > 0 {
		return &FieldError{Message: "重复的值模板列表", Field: "templates." + t.Templates[indexes[0]]}
	}

	if len(t.Templates) == 0 {
		t.Templates = []string{"post.xsl"}
	}

	for index, s := range t.Screenshots {
		if !utils.FileExists(filepath.Join(dir, s)) {
			return &FieldError{Message: "不存在的示例图", Field: "screenshots[" + strconv.Itoa(index) + "]"}
		}
	}
	indexes = sliceutil.Dup(t.Templates, func(i, j int) bool { return t.Screenshots[i] == t.Screenshots[j] })
	if len(indexes) > 0 {
		return &FieldError{Message: "重复的值示例图", Field: "screenshots[" + strconv.Itoa(indexes[0]) + "]"}
	}

	return nil
}

// LoadTheme 加载指定主题
func LoadTheme(dir, name string) (*Theme, error) {
	dir = filepath.Join(dir, "themes", name)
	path := filepath.Join(dir, "theme.yaml")

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
