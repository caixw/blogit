// SPDX-License-Identifier: MIT

package data

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2/internal/loader"
)

func TestRobots(t *testing.T) {
	a := assert.New(t, false)

	cfg := &loader.Config{
		Robots: []*loader.Agent{
			{
				Agent:    []string{"*"},
				Disallow: []string{"/themes/"},
			},
		},
	}

	robots := newRobots(cfg, nil)
	a.Empty(robots.Sitemaps).Equal(len(robots.Agents), 1)

	robots = newRobots(cfg, &Sitemap{Permalink: "/sitemap.xml"})
	a.Equal(robots.Sitemaps[0], "/sitemap.xml").Equal(len(robots.Agents), 1)
}
