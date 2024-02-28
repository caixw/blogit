// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
)

func (b *Builder) buildPosts(d *data.Data) error {
	for _, p := range d.Posts {
		page := b.page(p.Template)
		page.Title = p.Title + d.TitleSuffix
		page.Permalink = p.Permalink
		page.Keywords = p.Keywords
		page.Description = p.Summary
		page.Language = d.Language
		page.Post = p
		page.JSONLD = p.JSONLD
		page.License = p.License

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

	return nil
}

func (b *Builder) buildIndexes(d *data.Data) error {
	for _, index := range d.Indexes {
		page := b.page(vars.IndexTemplate)
		if index.Index == 1 {
			page.Title = d.Title
		} else {
			page.Title = index.Title + d.TitleSuffix
		}
		page.Permalink = index.Permalink
		page.Keywords = index.Keywords
		page.Description = index.Description
		page.Language = d.Language
		page.Index = index

		if index.Next != nil {
			page.Next = &loader.Link{
				URL:  index.Next.Permalink,
				Text: index.Next.Title,
			}
		}
		if index.Prev != nil {
			page.Prev = &loader.Link{
				URL:  index.Prev.Permalink,
				Text: index.Prev.Title,
			}
		}

		if err := b.appendTemplateFile(index.Path, page); err != nil {
			return err
		}
	}

	return nil
}
