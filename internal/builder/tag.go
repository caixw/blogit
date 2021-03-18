// SPDX-License-Identifier: MIT

package builder

import (
	"strings"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildTags(d *data.Data) error {
	keys := make([]string, 0, len(d.Tags.Tags))

	for _, t := range d.Tags.Tags {
		p := b.page(vars.TagTemplate)
		p.Title = t.Title
		p.Permalink = t.Permalink
		p.Keywords = t.Title + "," + t.Slug
		p.Description = t.Content
		p.Language = d.Language
		p.Tag = t

		if t.Next != nil {
			p.Next = &loader.Link{
				URL:  t.Next.Permalink,
				Text: t.Next.Title,
			}
		}
		if t.Prev != nil {
			p.Prev = &loader.Link{
				URL:  t.Prev.Permalink,
				Text: t.Prev.Title,
			}
		}

		if err := b.appendTemplateFile(t.Path, p); err != nil {
			return err
		}

		keys = append(keys, t.Slug, t.Title)
	}

	p := b.page(vars.TagTemplate)
	p.Permalink = d.Tags.Permalink
	p.Keywords = strings.Join(keys, ",")
	p.Language = d.Language
	p.Tags = d.Tags.Tags
	return b.appendTemplateFile(vars.TagsFilename, p)
}
