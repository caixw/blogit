// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func (b *builder) buildProfile(d *data.Data) error {
	if d.Profile == nil {
		return nil
	}

	buf := &errwrap.Buffer{}

	buf.Printf("<!-- 当前文件由 %s 自动生成，请勿手动修改 -->", vars.URL)
	buf.WByte('\n').WByte('\n')

	p := d.Profile
	if p.Alternate != nil {
		buf.WString(p.Alternate.Content)
	} else {
		buf.WString(p.Title).WByte('\n').WByte('\n')

		for _, p := range p.Posts {
			buf.Printf("- [%s](%s)\n", p.Title, p.Permalink)
		}
		buf.WByte('\n')

		buf.WString(p.Footer).WByte('\n')
	}

	if buf.Err != nil {
		return buf.Err
	}

	b.files[p.Path] = buf.Bytes()
	return nil
}
