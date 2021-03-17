// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildArchive(d *data.Data) error {
	p := b.page(vars.ArchiveTemplate)
	p.Title = d.Archives.Title
	p.Permalink = d.Archives.Permalink
	p.Keywords = "TODO"
	p.Description = "TODO"
	p.Language = d.Language

	return b.appendTemplateFile(vars.ArchiveFilename, p)
}
