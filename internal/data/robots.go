// SPDX-License-Identifier: MIT

package data

import "github.com/caixw/blogit/v2/internal/loader"

// Robots robots.txt 的相关内容
type Robots struct {
	Path     string
	Sitemaps []string
	Agents   []*loader.Agent
}

func newRobots(conf *loader.Config, sitemap *Sitemap) *Robots {
	robots := &Robots{
		Path:   "robots.txt",
		Agents: conf.Robots,
	}

	if sitemap != nil {
		robots.Sitemaps = []string{sitemap.Permalink}
	}

	return robots
}
