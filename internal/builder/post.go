// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildPosts(d *data.Data) error {
	for _, p := range d.Index.Posts {
		page := b.page(p.Template)
		page.Title = p.Title
		page.Permalink = p.Permalink
		page.Keywords = "todo"
		page.Description = p.Summary
		page.Language = d.Language
		page.Post = p

		if p.Next != nil {
			page.Next = &loader.Link{
				URL:  p.Next.Permalink,
				Text: p.Next.Title,
			}
		}
		if p.Prev != nil {
			page.Prev = &loader.Link{
				URL:  p.Prev.Permalink,
				Text: p.Prev.Title,
			}
		}

		if err := b.appendTemplateFile(p.Path, page); err != nil {
			return err
		}
	}

	page := b.page(vars.IndexTemplate)
	page.Permalink = d.URL
	page.Keywords = "todo"
	page.Description = "todo"
	page.Language = d.Language
	page.Posts = d.Index.Posts

	return b.appendTemplateFile(vars.IndexFilename, page)
}
