// SPDX-License-Identifier: MIT

package loader

// Profile 用于生成 github.com/profile 中的 README.md 内容
type Profile struct {
	Alternate string `yaml:"alternate"` // 采用此文件的内容代替

	// 当 alternate 为空，以下值才生效
	Title string `yaml:"title"`
	Size  uint   `yaml:"size"` // 显示的条数
}

func (p *Profile) sanitize() *FieldError {
	if p.Alternate != "" {
		switch {
		case p.Title != "":
			return &FieldError{Field: "title", Message: "只能为空"}
		case p.Size != 0:
			return &FieldError{Field: "size", Message: "只能为空"}
		}
		return nil
	}

	if p.Title == "" {
		return &FieldError{Field: "title", Message: "不能为空"}
	}

	if p.Size <= 0 {
		return &FieldError{Field: "size", Message: "必须大于 0"}
	}

	return nil
}
