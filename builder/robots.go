// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/vars"
)

func (b *Builder) buildRobots(d *data.Data) error {
	if d.Robots == nil {
		return nil
	}

	buf := &errwrap.Buffer{}

	buf.Printf("# %s", vars.FileHeader).WByte('\n').WByte('\n')

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

	return b.appendFile(d.Robots.Path, buf.Bytes())
}
