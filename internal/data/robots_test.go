// SPDX-License-Identifier: MIT

package data

import (
	"testing"

	"github.com/caixw/blogit/internal/loader"
	"github.com/issue9/assert"
)

func TestRobots(t *testing.T) {
	a := assert.New(t)

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
