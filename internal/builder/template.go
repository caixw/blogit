// SPDX-License-Identifier: MIT

package builder

import (
	"html/template"
	"path/filepath"
	"regexp"
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func newTemplate(d *data.Data, src string) (*template.Template, error) {
	templateFuncs := template.FuncMap{
		"strip":   stripTags,
		"html":    htmlEscaped,
		"rfc3339": rfc3339,
		"themeURL": func(p string) string {
			return data.BuildURL(d.URL, p)
		},
	}

	return template.New(d.Theme.ID).
		Funcs(templateFuncs).
		ParseGlob(filepath.Join(src, vars.ThemesDir, d.Theme.ID, vars.LayoutDir, "/*"))
}

func rfc3339(t time.Time) interface{} {
	return t.Format(time.RFC3339)
}

// 将内容显示为 HTML 内容
func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}
