// SPDX-License-Identifier: MIT

// Package vars 提供各种代码级别的配置项
package vars

// 各个配置项
const (
	ConfYAML  = "conf.yaml"
	TagsYAML  = "tags.yaml"
	ThemeYAML = "theme.yaml"

	ThemesDir = "themes"
	PostsDir  = "posts"
	TagsDir   = "tags"

	// DefaultTemplate 默认的模板名称
	//
	// 在文章未指定模板时，都将采用此模板作为其转换方式。
	DefaultTemplate = "post.xsl"
)
