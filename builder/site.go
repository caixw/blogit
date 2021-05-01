// SPDX-License-Identifier: MIT

package builder

import (
	"bytes"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

type page struct {
	Site *site

	// 当前页面类型，同时也是当前页的模板名称。
	// 除了文章详情页以外，其它值是固定的，
	// 文章详情页可以指定不同的模板，可以不同。
	Type string

	Title       string // 标题，html>head>title 的内容，会带上后缀。
	Permalink   string // 当前页的唯一链接
	Keywords    string
	Description string
	Prev        *loader.Link
	Next        *loader.Link
	Authors     []*loader.Author
	License     *loader.Link
	Language    string
	JSONLD      string // JSON-LD 数据

	// 以下内容，仅在对应的页面才会有内容
	Tag      *data.Tag   // 标签详细页面，非标签详细页，则为空
	Post     *data.Post  // 文章详细内容，仅文章页面用到。
	Index    *data.Index // 索引页的数据
	Archives *data.Archives
}

type site struct {
	AppName    string // 程序名称
	AppURL     string // 程序官网
	AppVersion string // 当前程序的版本号
	Theme      *data.Theme
	Highlights []*styleLink // 高亮引用的 CSS 文件列表

	Title    string
	Subtitle string       // 网站副标题
	URL      string       // 网站地址，若是一个子目录，则需要包含该子目录
	Icon     *loader.Icon // 网站图标
	Author   *loader.Author
	RSS      *loader.Link // RSS 指针方便模板判断其值是否为空
	Atom     *loader.Link
	Sitemap  *loader.Link
	Menus    []*loader.Link
	Tags     *data.Tags

	Uptime   time.Time
	Created  time.Time
	Modified time.Time
	Builded  time.Time // 最后次编译时间
}

type styleLink struct {
	Media string
	URL   string
}

func newSite(d *data.Data) *site {
	s := &site{
		AppName:    vars.Name,
		AppURL:     vars.URL,
		AppVersion: vars.Version(),
		Theme:      d.Theme,
		Highlights: make([]*styleLink, 0, len(d.Highlights)),
		Menus:      d.Menus,

		Title:    d.Title,
		Subtitle: d.Subtitle,
		URL:      d.URL,
		Icon:     d.Icon,
		Author:   d.Author,
		Tags:     d.Tags,

		Uptime:   d.Uptime,
		Created:  d.Created,
		Modified: d.Modified,
		Builded:  d.Builded,
	}

	if d.RSS != nil {
		s.RSS = &loader.Link{URL: d.RSS.Permalink, Text: d.RSS.Title}
	}
	if d.Atom != nil {
		s.Atom = &loader.Link{URL: d.Atom.Permalink, Text: d.Atom.Title}
	}
	if d.Sitemap != nil {
		s.Sitemap = &loader.Link{URL: d.Sitemap.Permalink, Text: d.Sitemap.Title}
	}

	for _, h := range d.Highlights {
		s.Highlights = append(s.Highlights, &styleLink{
			Media: h.Media,
			URL:   h.URL,
		})
	}

	return s
}

func (b *Builder) page(t string) *page {
	return &page{
		Site: b.site,
		Type: t,
	}
}

var cssFormatter = html.New(
	html.ClassPrefix(vars.HighlightClassPrefix),
)

func (b *Builder) buildHighlights(d *data.Data) error {
	for _, h := range d.Highlights {
		buf := &bytes.Buffer{}
		if err := cssFormatter.WriteCSS(buf, styles.Get(h.Name)); err != nil {
			return err
		}
		if err := b.appendFile(h.Path, buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
