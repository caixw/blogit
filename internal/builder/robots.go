// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

func (b *Builder) buildRobots(d *data.Data) error {
	if d.Robots == nil {
		return nil
	}

	buf := &errwrap.Buffer{}

	buf.Printf("# 当前文件由 %s 自动生成，请勿手动修改", vars.URL)
	buf.WByte('\n').WByte('\n')

	for _, agent := range d.Robots.Agents {
		for _, a := range agent.Agent {
			buf.Println("User-agent:", a)
		}

		for _, disallow := range agent.Disallow {
			buf.Println("Disallow:", disallow)
		}

		for _, allow := range agent.Allow {
			buf.Println("Allow:", allow)
		}

		buf.WByte('\n')
	}

	for _, sitemap := range d.Robots.Sitemaps {
		buf.Println("Sitemap:", sitemap)
	}

	if buf.Err != nil {
		return buf.Err
	}

	b.appendFile(d.Robots.Path, "", buf.Bytes())
	return nil
}
