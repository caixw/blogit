// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/vars"
)

func (b *Builder) buildArchive(d *data.Data) error {
	p := b.page(vars.ArchiveTemplate)
	p.Title = d.Archives.Title + d.TitleSuffix
	p.Permalink = d.Archives.Permalink
	p.Keywords = d.Archives.Keywords
	p.Description = d.Archives.Description
	p.Language = d.Language
	p.Archives = d.Archives

	return b.appendTemplateFile(vars.ArchiveFilename, p)
}
