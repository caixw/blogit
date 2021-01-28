// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

type tags struct {
	Info *info
	Tags []*tag
}

type tag struct {
	Info      *info
	Permalink string
	Title     string
	Created   time.Time
	Modified  time.Time
	Posts     []*postMeta
	Content   string
}

func newTag(t *data.Tag, d *data.Data, i *info) *tag {
	ps := make([]*postMeta, 0, len(t.Posts))
	for _, p := range t.Posts {
		ps = append(ps, &postMeta{
			Permalink: d.BuildURL(p.Path),
			Language:  p.Language,
			Title:     p.Title,
			Created:   p.Created,
			Modified:  p.Modified,
			Summary:   p.Summary,
		})
	}

	return &tag{
		Info:      i,
		Permalink: d.BuildURL(t.Path),
		Title:     t.Title,
		Created:   t.Created,
		Modified:  t.Modified,
		Posts:     ps,
		Content:   t.Content,
	}
}

func (b *builder) buildTags(d *data.Data, i *info) error {
	ts := make([]*tag, 0, len(d.Tags))

	for _, t := range d.Tags {
		tt := newTag(t, d, i)
		if err := b.appendTemplateFile(d, t.Path, vars.TagTemplate, tt); err != nil {
			return err
		}

		tt.Posts = nil
		ts = append(ts, tt)
	}

	t := &tags{
		Info: i,
		Tags: ts,
	}
	return b.appendTemplateFile(d, vars.TagsFilename, vars.TagsTemplate, t)
}
