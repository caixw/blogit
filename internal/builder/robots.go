// SPDX-License-Identifier: MIT

package builder

import (
	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
)

func (b *builder) buildRobots(d *data.Data) error {
	buf := &errwrap.Buffer{}

	for _, agent := range d.Robots.Agents {
		buf.Printf("User-agent: %s\n", agent.Agent)

		for _, disallow := range agent.Disallow {
			buf.Printf("Disallow: %s\n", disallow)
		}

		for _, allow := range agent.Allow {
			buf.Printf("Disallow: %s\n", allow)
		}

		buf.WByte('\n')
	}

	if d.Robots.Sitemap != "" {
		buf.Println("Sitemap:", d.Robots.Sitemap)
	}

	if buf.Err != nil {
		return buf.Err
	}

	b.files[d.Robots.Path] = buf.Bytes()
	return nil
}
