// SPDX-License-Identifier: MIT

package builder

import "github.com/caixw/blogit/internal/data"

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

	Atom    bool `xml:"atom,omitempty"`
	RSS     bool `xml:"rss,omitempty"`
	Sitemap bool `xml:"sitemap,omitempty"`
	Archive bool `xml:"archive,omitempty"`

	Uptime   string `xml:"uptime,attr"`
	Created  string `xml:"created,attr"`
	Modified string `xml:"modified,attr"`
	Builded  string `xml:"builded,attr"` // 最后次编译时间
}

type menu struct {
	// 链接对应的图标。可以是字体图标或是图片链接，模板根据情况自动选择。
	Icon string `xml:"icon,attr"`
	URL  string `xml:"url,attr"`  // 链接地址
	Text string `xml:",chardata"` // 链接的文本
}

type icon struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"` // mime type
}

func (b *Builder) buildInfo(path string, d *data.Data) error {
	menus := make([]*menu, 0, len(d.Menus))
	for _, m := range d.Menus {
		menus = append(menus, &menu{
			Icon: m.Icon,
			URL:  m.URL,
			Text: m.Text,
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
			URL:  d.License.URL,
			Text: d.License.Text,
		},

		Uptime:   ft(d.Uptime),
		Created:  ft(d.Created),
		Modified: ft(d.Modified),
		Builded:  ft(d.Builded),
	}

	i.Atom = d.Atom != nil
	i.RSS = d.RSS != nil
	i.Sitemap = d.Sitemap != nil
	i.Archive = d.Archive != nil

	return b.appendXMLFile(d, path, "", d.Modified, i)
}
