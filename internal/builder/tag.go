// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildTags(d *data.Data) error {
	for _, t := range d.Tags.Tags {
		p := b.page(vars.TagTemplate)
		p.Title = t.Title
		p.Permalink = t.Permalink
		p.Keywords = t.Keywords
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
	}

	p := b.page(vars.TagTemplate)
	p.Permalink = d.Tags.Permalink
	p.Keywords = d.Tags.Keywords
	p.Description = d.Tags.Description
	p.Language = d.Language
	p.Tags = d.Tags.Tags
	return b.appendTemplateFile(vars.TagsFilename, p)
}
