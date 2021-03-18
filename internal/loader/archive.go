// SPDX-License-Identifier: MIT

package loader

// 归档的类型
const (
	ArchiveTypeYear  = "year"
	ArchiveTypeMonth = "month"
)

var defaultArchiveFormats = map[string]string{
	ArchiveTypeYear:  "2006",
	ArchiveTypeMonth: "2006-01",
}

// Archive 存档页的配置内容
type Archive struct {
	Title       string `yaml:"title"` // 存档页的标题
	Keywords    string `yaml:"keywords,omitempty"`
	Description string `yaml:"description,omitempty"`
	Order       string `yaml:"order,omitempty"`  // 排序方式
	Type        string `yaml:"type,omitempty"`   // 存档的分类方式，可以按年或是按月
	Format      string `yaml:"format,omitempty"` // 标题的格式化字符串，被 time.Format 所格式化。
}

func (a *Archive) sanitize() *FieldError {
	if a.Title == "" {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	if a.Type == "" {
		a.Type = ArchiveTypeYear
	} else {
		if a.Type != ArchiveTypeMonth && a.Type != ArchiveTypeYear {
			return &FieldError{Message: "取值不正确", Field: "type", Value: a.Type}
		}
	}

	if a.Order == "" {
		a.Order = OrderDesc
	} else {
		if a.Order != OrderAsc && a.Order != OrderDesc {
			return &FieldError{Message: "取值不正确", Field: "order", Value: a.Order}
		}
	}

	if a.Format == "" {
		a.Format = defaultArchiveFormats[a.Type]
		if a.Format == "" {
			panic(&FieldError{Message: "取值不正确", Field: "type"})
		}
	}

	return nil
}
