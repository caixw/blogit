// SPDX-License-Identifier: MIT

package loader

import (
	"path/filepath"
	"strconv"
)

// Theme 主题
type Theme struct {
	ID          string    `yaml:"-"`
	Description string    `yaml:"description,omitempty"`
	Authors     []*Author `yaml:"authors,omitempty"`
	Templates   []string  `yaml:"templates"`
	Screenshot  []string  `yaml:"screenshot,omitempty"`
}

func (t *Theme) sanitize() *FieldError {
	for index, author := range t.Authors {
		if err := author.sanitize(); err != nil {
			err.Field = "authors[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}
	return nil
}

// LoadTheme 加载指定主题
func LoadTheme(dir, name string) (*Theme, error) {
	path := filepath.Join(dir, "themes", name, "theme.yaml")

	theme := &Theme{}
	if err := loadYAML(path, &theme); err != nil {
		return nil, err
	}
	theme.ID = name

	if err := theme.sanitize(); err != nil {
		err.File = path
		return nil, err
	}
	return theme, nil
}
