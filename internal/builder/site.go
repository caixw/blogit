// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

type page struct {
	Site *site

	Title       string       // 标题
	Permalink   string       // 当前页的唯一链接
	Keywords    string       // meta.keywords 的值
	Description string       // meta.description 的值
	Prev        *loader.Link // 前一页
	Next        *loader.Link // 下一页
	Type        string       // 当前页面类型
	Authors     []*loader.Author
	License     *loader.Link // 当前页的版本信息，可以为空
	Language    string       // 页面语言

	// 以下内容，仅在对应的页面才会有内容
	Tag     *data.Tag      // 标签详细页面，非标签详细页，则为空
	Tags    []*data.Tag    // 标签列表页面，否则为空
	Posts   []*data.Post   // 文章列表，仅标签详情页和搜索页用到。
	Post    *data.Post     // 文章详细内容，仅文章页面用到。
	Archive *data.Archives // 归档
}

type site struct {
	AppName    string // 程序名称
	AppURL     string // 程序官网
	AppVersion string // 当前程序的版本号
	Theme      *loader.Theme

	TitleSuffix string // 非首页标题的后缀
	Title       string
	Subtitle    string       // 网站副标题
	URL         string       // 网站地址，若是一个子目录，则需要包含该子目录
	Icon        *loader.Icon // 网站图标
	RSS         *loader.Link // RSS 指针方便模板判断其值是否为空
	Atom        *loader.Link
	Sitemap     *loader.Link

	Uptime   time.Time
	Created  time.Time
	Modified time.Time
	Builded  time.Time // 最后次编译时间
}

func newSite(d *data.Data) *site {
	return &site{
		AppName:    vars.Name,
		AppURL:     vars.URL,
		AppVersion: vars.Version(),
		Theme:      d.Theme,

		TitleSuffix: d.TitleSuffix,
		Title:       d.Title,
		Subtitle:    d.Subtitle,
		URL:         d.URL,
		Icon:        d.Icon,
		RSS:         &loader.Link{URL: d.BuildURL(vars.RssXML), Text: d.RSS.Title},
		Atom:        &loader.Link{URL: d.BuildURL(vars.AtomXML), Text: d.Atom.Title},
		Sitemap:     &loader.Link{URL: d.BuildURL(vars.SitemapXML), Text: d.Title},

		Uptime:   d.Uptime,
		Created:  d.Created,
		Modified: d.Modified,
		Builded:  d.Builded,
	}
}

func (s *site) page(t string) *page {
	return &page{
		Site: s,
		Type: t,
	}
}

func (b *builder) page(t string) *page {
	return b.site.page(t)
}
