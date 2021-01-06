// SPDX-License-Identifier: MIT

package builder

import "github.com/caixw/blogit/internal/data"

type posts struct {
	XMLName struct{} `xml:"posts"`

	Posts []*postMeta `xml:"post"`
}

type postMeta struct {
	Permalink string   `xml:"permalink"`
	Title     string   `xml:"title"`
	Created   datetime `xml:"created"`
	Modified  datetime `xml:"modified"`
	Tags      []*tag   `xml:"tag,omitempty"`
}

type post struct {
	XMLName struct{} `xml:"post"`

	Permalink string    `xml:"permalink"`
	Title     string    `xml:"title"`
	Created   datetime  `xml:"created"`
	Modified  datetime  `xml:"modified"`
	Tags      []*tag    `xml:"tag"`
	Language  string    `xml:"language,attr"`
	Outdated  *outdated `xml:"outdated"`
	Authors   []*author `xml:"author"`
	License   *link     `xml:"license"`
	Content   string    `xml:"content"`
	Prev      *link     `xml:"prev"`
	Next      *link     `xml:"next"`
}

type author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

type link struct {
	URL   string `xml:"url"`
	Title string `xml:"title,attr,omitempty"`
	Text  string `xml:"text"`
}

type outdated struct {
	Outdated datetime `xml:"outdated"`
	Content  string   `xml:",chardata"`
}

func (b *Builder) buildPosts(d *data.Data) error {
	index := &posts{Posts: make([]*postMeta, 0, len(d.Posts))}

	for _, p := range d.Posts {
		tags := make([]*tag, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, &tag{
				Permalink: d.BuildURL("tags", t.Slug+".xml"),
				Title:     t.Title,
				Color:     t.Color,
				Content:   t.Content,
				Created:   toDatetime(t.Created, d),
				Modified:  toDatetime(t.Modified, d),
			})
		}

		authors := make([]*author, 0, len(p.Authors))
		for _, a := range p.Authors {
			authors = append(authors, &author{
				Name:   a.Name,
				URL:    a.URL,
				Email:  a.Email,
				Avatar: a.Avatar,
			})
		}

		pp := &post{
			Permalink: d.BuildURL(p.Slug + ".xml"),
			Title:     p.Title,
			Created:   toDatetime(p.Created, d),
			Modified:  toDatetime(p.Modified, d),
			Tags:      tags,
			Language:  p.Language,
			Outdated: &outdated{
				Outdated: toDatetime(p.Outdated.Outdated, d),
				Content:  p.Outdated.Content,
			},
			Authors: authors,
			License: &link{
				URL:   p.License.URL,
				Title: p.License.Title,
				Text:  p.License.Text,
			},
			Content: p.Content,
		}
		if p.Prev != nil {
			pp.Prev = &link{
				URL:   d.BuildURL(p.Prev.Slug + ".xml"),
				Title: "上一篇文章",
				Text:  p.Prev.Title,
			}
		}
		if p.Next != nil {
			pp.Next = &link{
				URL:   d.BuildURL(p.Next.Slug + ".xml"),
				Title: "上一篇文章",
				Text:  p.Next.Title,
			}
		}
		err := b.appendXMLFile(p.Slug+".xml", d.BuildThemeURL(p.Template), p.Modified, pp)
		if err != nil {
			return err
		}

		index.Posts = append(index.Posts, &postMeta{
			Permalink: d.BuildURL(p.Slug + ".xml"),
			Title:     p.Title,
			Created:   toDatetime(p.Created, d),
			Modified:  toDatetime(p.Modified, d),
			Tags:      tags,
		})
	}

	return b.appendXMLFile("index.xml", d.BuildThemeURL("index.xsl"), d.Modified, index)
}
