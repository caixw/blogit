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
		"html":    func(html string) interface{} { return template.HTML(html) },
		"js":      func(js string) interface{} { return template.JS(js) },
		"rfc3339": func(t time.Time) string { return t.Format(time.RFC3339) },
		"date":    func(t time.Time, format string) string { return t.Format(format) },
		"themeURL": func(p string) string {
			return data.BuildURL(d.URL, vars.ThemesDir, p)
		},
	}

	return template.New(d.Theme.ID).
		Funcs(templateFuncs).
		ParseGlob(filepath.Join(src, vars.ThemesDir, d.Theme.ID, vars.LayoutDir, "/*"))
}

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}
