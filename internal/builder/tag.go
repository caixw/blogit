// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

type tags struct {
	XMLName struct{} `xml:"tags"`
	Tags    []*tag   `xml:"tag"`
}

type tag struct {
	XMLName   struct{}    `xml:"tag"`
	Permalink string      `xml:"permalink,attr"`
	Title     string      `xml:"title"`
	Created   string      `xml:"created,attr,omitempty"`
	Modified  string      `xml:"modified,attr,omitempty"`
	Posts     []*postMeta `xml:"post,omitempty"`
	Content   *innerhtml  `xml:"summary,omitempty"`
}

func newTag(t *data.Tag, d *data.Data) *tag {
	ps := make([]*postMeta, 0, len(t.Posts))
	for _, p := range t.Posts {
		ps = append(ps, &postMeta{
			Permalink: d.BuildURL(p.Slug),
			Title:     p.Title,
			Created:   ft(p.Created),
			Modified:  ft(p.Modified),
			Summary:   newHTML(p.Summary),
		})
	}

	return &tag{
		Permalink: d.BuildURL(t.Path),
		Title:     t.Title,
		Created:   ft(t.Created),
		Modified:  ft(t.Modified),
		Posts:     ps,
		Content:   newHTML(t.Content),
	}
}

func (b *builder) buildTags(d *data.Data) error {
	tags := make([]*tag, 0, len(d.Tags))

	for _, t := range d.Tags {
		tt := newTag(t, d)
		if err := b.appendXMLFile(d, t.Path, d.Theme.Tag, t.Modified, tt); err != nil {
			return err
		}

		tt.Posts = nil
		tags = append(tags, tt)
	}

	return b.appendXMLFile(d, vars.TagsXML, d.Theme.Tags, d.Modified, tags)
}
