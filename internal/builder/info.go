// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
)

type info struct {
	XMLName struct{} `xml:"info"`

	URL         string    `xml:"url,attr"`
	Title       string    `xml:"title"`
	Subtitle    string    `xml:"subtitle"`
	TitleSuffix string    `xml:"titleSuffix"` // 每篇文章标题的后缀
	Icon        *icon     `xml:"icon"`
	Menus       []*menu   `xml:"menus"`
	Language    string    `xml:"language,attr"`
	Authors     []*author `xml:"author"`
	License     *link     `xml:"license"`

	Atom    *link `xml:"atom,omitempty"`
	RSS     *link `xml:"rss,omitempty"`
	Sitemap *link `xml:"sitemap,omitempty"`

	Uptime   string `xml:"uptime,attr"`
	Created  string `xml:"created,attr"`
	Modified string `xml:"modified,attr"`
	Builded  string `xml:"builded,attr"` // 最后次编译时间
}

type menu struct {
	// 链接对应的图标。可以是字体图标或是图片链接，模板根据情况自动选择。
	Icon  string `xml:"icon"`
	Title string `xml:"title"` // 链接的 title 属性
	URL   string `xml:"url"`   // 链接地址
	Text  string `xml:"text"`  // 链接的文本
}

type icon struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"` // mime type
}

func (b *Builder) buildInfo(path string, d *data.Data) error {
	menus := make([]*menu, 0, len(d.Menus))
	for _, m := range d.Menus {
		menus = append(menus, &menu{
			Icon:  m.Icon,
			Title: m.Title,
			URL:   m.URL,
			Text:  m.Text,
		})
	}

	authors := make([]*author, 0, len(d.Authors))
	for _, a := range d.Authors {
		authors = append(authors, &author{
			Name:   a.Name,
			URL:    a.URL,
			Email:  a.Email,
			Avatar: a.Avatar,
		})
	}

	i := info{
		URL:         d.URL,
		Title:       d.Title,
		Subtitle:    d.Subtitle,
		TitleSuffix: d.TitleSuffix,
		Icon:        &icon{URL: d.Icon.URL, Type: d.Icon.Type},
		Menus:       menus,
		Language:    d.Language,
		Authors:     authors,
		License: &link{
			URL:   d.License.URL,
			Title: d.License.Title,
			Text:  d.License.Text,
		},

		Uptime:   ft(d.Uptime),
		Created:  ft(d.Created),
		Modified: ft(d.Modified),
		Builded:  ft(d.Builded),
	}

	if d.Atom != nil {
		i.Atom = &link{
			URL:   d.BuildURL("atom.xml"),
			Title: d.Atom.Title,
			Text:  d.Atom.Title,
		}
	}

	if d.RSS != nil {
		i.RSS = &link{
			URL:   d.BuildURL("rss.xml"),
			Title: d.RSS.Title,
			Text:  d.RSS.Title,
		}
	}

	if d.Sitemap != nil {
		i.Sitemap = &link{
			URL:   d.BuildURL("sitemap.xml"),
			Title: "sitemap",
			Text:  "sitemap",
		}
	}

	return b.appendXMLFile(path, "", d.Modified, i)
}
