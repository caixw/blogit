// SPDX-License-Identifier: MIT

package loader

import (
	"path/filepath"
	"strconv"

	"github.com/caixw/blogit/internal/utils"
)

// Theme 主题
type Theme struct {
	ID          string    `yaml:"-"`
	Dir         string    `yaml:"-"` // 当前主题所在的目录
	Description string    `yaml:"description,omitempty"`
	Authors     []*Author `yaml:"authors,omitempty"`
	Templates   []string  `yaml:"templates"`
	Screenshots []string  `yaml:"screenshots,omitempty"`
}

func (t *Theme) sanitize() *FieldError {
	for index, author := range t.Authors {
		if err := author.sanitize(); err != nil {
			err.Field = "authors[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	for index, tpl := range t.Templates {
		if !utils.FileExists(filepath.Join(t.Dir, tpl)) {
			return &FieldError{Message: "不存在", Field: "templates[" + strconv.Itoa(index) + "]"}
		}
	}

	for index, s := range t.Screenshots {
		if !utils.FileExists(filepath.Join(t.Dir, s)) {
			return &FieldError{Message: "不存在", Field: "screenshots[" + strconv.Itoa(index) + "]"}
		}
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
	theme.ID = name
	theme.Dir = dir

	if err := theme.sanitize(); err != nil {
		err.File = path
		return nil, err
	}
	return theme, nil
}
