// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
)

type info struct {
	URL         string
	Title       string
	Subtitle    string
	TitleSuffix string // 每篇文章标题的后缀
	Icon        *icon
	Menus       []*menu
	Language    string
	Authors     []*author
	License     *link
	Theme       *theme

	Atom    bool
	RSS     bool
	Sitemap bool
	Archive bool

	Uptime   time.Time
	Created  time.Time
	Modified time.Time
	Builded  time.Time // 最后次编译时间
}

type theme struct {
	ID          string
	Base        string
	Description string
	Authors     []*author
}

type menu struct {
	// 链接对应的图标。可以是字体图标或是图片链接，模板根据情况自动选择。
	Icon string
	URL  string // 链接地址
	Text string // 链接的文本
}

type icon struct {
	URL  string
	Type string // mime type
}

func (b *builder) buildInfo(d *data.Data) *info {
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

	themeAuthors := make([]*author, 0, len(d.Theme.Authors))
	for _, a := range d.Theme.Authors {
		themeAuthors = append(themeAuthors, &author{
			Name:   a.Name,
			URL:    a.URL,
			Email:  a.Email,
			Avatar: a.Avatar,
		})
	}

	i := &info{
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
		Theme: &theme{
			ID:          d.Theme.ID,
			Base:        d.BuildThemeURL(),
			Description: d.Theme.Description,
			Authors:     themeAuthors,
		},

		Uptime:   d.Uptime,
		Created:  d.Created,
		Modified: d.Modified,
		Builded:  d.Builded,
	}

	i.Atom = d.Atom != nil
	i.RSS = d.RSS != nil
	i.Sitemap = d.Sitemap != nil
	i.Archive = d.Archive != nil

	return i
}
