// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

type posts struct {
	XMLName struct{}    `xml:"posts"`
	Posts   []*postMeta `xml:"post"`
}

type postMeta struct {
	Permalink string     `xml:"permalink,attr"`
	Title     string     `xml:"title"`
	Created   string     `xml:"created,attr,omitempty"`
	Modified  string     `xml:"modified,attr,omitempty"`
	Tags      []*tagMeta `xml:"tag,omitempty"`
	Summary   *innerhtml `xml:"summary,omitempty"`
}

type tagMeta struct {
	Permalink string `xml:"permalink,attr"`
	Title     string `xml:",chardata"`
}

type post struct {
	XMLName   struct{}   `xml:"post"`
	Permalink string     `xml:"permalink,attr"`
	Title     string     `xml:"title"`
	Created   string     `xml:"created,attr,omitempty"`
	Modified  string     `xml:"modified,attr,omitempty"`
	Tags      []*tagMeta `xml:"tag"`
	Language  string     `xml:"language,attr"`
	Outdated  *outdated  `xml:"outdated,omitempty"`
	Authors   []*author  `xml:"author"`
	License   *link      `xml:"license"`
	Summary   *innerhtml `xml:"summary,omitempty"`
	Content   *innerhtml `xml:"content"`
	Prev      *link      `xml:"prev"`
	Next      *link      `xml:"next"`
}

type author struct {
	Name   string `xml:",chardata"`
	URL    string `xml:"url,attr,omitempty"`
	Email  string `xml:"email,attr,omitempty"`
	Avatar string `xml:"avatar,attr,omitempty"`
}

type link struct {
	URL  string `xml:"url,attr"`
	Text string `xml:",chardata"`
}

type outdated struct {
	Outdated string `xml:"outdated,attr,omitempty"` // 过期的时间
	Content  string `xml:",chardata"`
}

func (b *Builder) buildPosts(d *data.Data) error {
	index := &posts{Posts: make([]*postMeta, 0, len(d.Posts))}

	for _, p := range d.Posts {
		tags := make([]*tagMeta, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, &tagMeta{
				Permalink: d.BuildURL(t.Path),
				Title:     t.Title,
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

		var od *outdated
		if p.Outdated != nil {
			od = &outdated{
				Content:  p.Outdated.Content,
				Outdated: ft(p.Outdated.Outdated),
			}
		}

		pp := &post{
			Permalink: d.BuildURL(p.Path),
			Title:     p.Title,
			Created:   ft(p.Created),
			Modified:  ft(p.Modified),
			Tags:      tags,
			Language:  p.Language,
			Outdated:  od,
			Authors:   authors,
			License: &link{
				URL:  p.License.URL,
				Text: p.License.Text,
			},
			Content: newHTML(p.Content),
			Summary: newHTML(p.Summary),
		}
		if p.Prev != nil {
			pp.Prev = &link{
				URL:  d.BuildURL(p.Prev.Path),
				Text: p.Prev.Title,
			}
		}
		if p.Next != nil {
			pp.Next = &link{
				URL:  d.BuildURL(p.Next.Path),
				Text: p.Next.Title,
			}
		}
		err := b.appendXMLFile(d, p.Path, p.Template, p.Modified, pp)
		if err != nil {
			return err
		}

		index.Posts = append(index.Posts, &postMeta{
			Permalink: d.BuildURL(p.Path),
			Title:     p.Title,
			Created:   ft(p.Created),
			Modified:  ft(p.Modified),
			Tags:      tags,
			Summary:   newHTML(p.Summary),
		})
	}

	return b.appendXMLFile(d, vars.IndexXML, d.Theme.Index, d.Modified, index)
}
