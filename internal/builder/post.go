// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildPosts(d *data.Data) error {
	for _, p := range d.Posts {
		page := b.page(p.Template)
		page.Title = p.Title
		page.Permalink = p.Permalink
		page.Keywords = "todo"
		page.Description = p.Summary
		page.Language = d.Language
		page.Post = p
		// TODO prev

		if err := b.appendTemplateFile(p.Path, page); err != nil {
			return err
		}
	}

	page := b.page(vars.IndexTemplate)
	page.Permalink = d.URL
	page.Keywords = "todo"
	page.Description = "todo"
	page.Language = d.Language
	page.Posts = d.Posts

	return b.appendTemplateFile(vars.IndexFilename, page)
}
