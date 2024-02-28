// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/vars"
)

func (b *Builder) buildProfile(d *data.Data) error {
	if d.Profile == nil {
		return nil
	}
	p := d.Profile

	buf := &errwrap.Buffer{}

	buf.Printf("<!-- %s -->", vars.FileHeader).WByte('\n').WByte('\n')

	buf.WString(p.Title).WByte('\n').WByte('\n')

	for _, post := range p.Posts {
		buf.Printf("- [%s](%s)\n", post.Title, post.Permalink)
	}
	buf.WByte('\n')

	buf.WString(p.Footer).WByte('\n')

	if buf.Err != nil {
		return buf.Err
	}

	return b.appendFile(p.Path, buf.Bytes())
}
