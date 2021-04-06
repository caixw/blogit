// SPDX-License-Identifier: MIT

// Package vars 提供各种代码级别的配置项
package vars

// 各个配置项
const (
	Name = "blogit"
	URL  = "https://github.com/caixw/blogit"

	ConfYAML  = "conf.yaml"
	TagsYAML  = "tags.yaml"
	ThemeYAML = "theme.yaml"

	ThemesDir = "themes"
	PostsDir  = "posts"
	TagsDir   = "tags"
	LayoutDir = "layout"

	TagsFilename        = "tags" + Ext
	IndexFilename       = "index" + Ext    // 首页
	IndexFilenameFormat = "index-%d" + Ext // 非首页的索引页
	ArchiveFilename     = "archive" + Ext
	RssXML              = "rss.xml"
	AtomXML             = "atom.xml"
	SitemapXML          = "sitemap.xml"

	DefaultTemplate = "post"
	IndexTemplate   = "index"
	TagTemplate     = "tag"
	TagsTemplate    = "tags"
	ArchiveTemplate = "archive"

	Ext         = ".html" // 生成后的文件后缀名
	MarkdownExt = ".md"

	// 默认的高亮主题色
	// 值可以从 github.com/alecthomas/chroma/styles 获取
	HighlightClassPrefix = "hl-" // 语法高亮的统一类名前缀
)
