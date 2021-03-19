// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildArchive(d *data.Data) error {
	p := b.page(vars.ArchiveTemplate)
	p.Title = d.Archives.Title + d.TitleSuffix
	p.Permalink = d.Archives.Permalink
	p.Keywords = d.Archives.Keywords
	p.Description = d.Archives.Description
	p.Language = d.Language
	p.Archives = d.Archives.Archives

	return b.appendTemplateFile(vars.ArchiveFilename, p)
}
