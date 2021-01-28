// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

type posts struct {
	Info  *info
	Posts []*postMeta
}

type postMeta struct {
	Permalink string
	Title     string
	Language  string
	Created   time.Time
	Modified  time.Time
	Tags      []*tagMeta
	Summary   string
}

type tagMeta struct {
	Permalink string
	Title     string
}

type post struct {
	Info      *info
	Permalink string
	Title     string
	Created   time.Time
	Modified  time.Time
	Tags      []*tagMeta
	Language  string
	Outdated  *outdated
	Authors   []*author
	License   *link
	Summary   string
	Content   string
	Image     string
	Prev      *link
	Next      *link
}

type author struct {
	Name   string
	URL    string
	Email  string
	Avatar string
}

type link struct {
	URL  string
	Text string
}

type outdated struct {
	Outdated time.Time // 过期的时间
	Content  string
}

func (b *builder) buildPosts(d *data.Data, i *info) error {
	index := &posts{
		Info:  i,
		Posts: make([]*postMeta, 0, len(d.Posts)),
	}

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
				Outdated: p.Outdated.Outdated,
			}
		}

		pp := &post{
			Info:      i,
			Permalink: d.BuildURL(p.Path),
			Title:     p.Title,
			Created:   p.Created,
			Modified:  p.Modified,
			Tags:      tags,
			Language:  p.Language,
			Outdated:  od,
			Authors:   authors,
			License: &link{
				URL:  p.License.URL,
				Text: p.License.Text,
			},
			Content: p.Content,
			Summary: p.Summary,
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
		err := b.appendTemplateFile(d, p.Path, p.Template, pp)
		if err != nil {
			return err
		}

		index.Posts = append(index.Posts, &postMeta{
			Permalink: d.BuildURL(p.Path),
			Language:  p.Language,
			Title:     p.Title,
			Created:   p.Created,
			Modified:  p.Modified,
			Tags:      tags,
			Summary:   p.Summary,
		})
	}

	return b.appendTemplateFile(d, vars.IndexFilename, vars.IndexTemplate, index)
}
