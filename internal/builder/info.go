// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
)

type info struct {
	URL      string
	Title    string
	Subtitle string
	Icon     *icon
	Language string
	Authors  []*author
	License  *link
	Theme    *theme

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

type icon struct {
	URL  string
	Type string // mime type
}

func (b *builder) buildInfo(d *data.Data) *info {
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
		URL:      d.URL,
		Title:    d.Title,
		Subtitle: d.Subtitle,
		Icon:     &icon{URL: d.Icon.URL, Type: d.Icon.Type},
		Language: d.Language,
		Authors:  authors,
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
