// SPDX-License-Identifier: MIT

package builder

import (
	"strconv"
	"time"

	"github.com/caixw/blogit/internal/data"
)

const sitempaNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	XMLName struct{} `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLSet  []*url   `xml:"url,omitempty"`
}

type url struct {
	Loc        string `xml:"loc"`
	Lastmod    string `xml:"lastmod"`
	Changefreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

func (b *Builder) buildSitemap(path string, d *data.Data) error {
	s := &urlset{
		XMLNS:  sitempaNamespace,
		URLSet: make([]*url, 0, len(d.Tags)+len(d.Posts)),
	}

	conf := d.Sitemap
	if conf.EnableTag {
		s.append(d.BuildURL("tags.xml"), d.Modified, conf.Changefreq, conf.Priority)
		for _, tag := range d.Tags {
			s.append(d.BuildURL(tag.Slug+".xml"), tag.Modified, conf.Changefreq, conf.Priority)
		}
	}

	s.append(d.URL, d.Modified, conf.Changefreq, conf.Priority)
	for _, p := range d.Posts {
		s.append(d.BuildURL(p.Slug+".xml"), p.Modified, conf.PostChangefreq, conf.PostPriority)
	}

	return b.appendXMLFile(path, conf.XSL, d.Modified, s)
}

func (us *urlset) append(loc string, lastmod time.Time, changefreq string, priority float64) {
	us.URLSet = append(us.URLSet, &url{
		Loc:        loc,
		Lastmod:    lastmod.Format(time.RFC3339),
		Changefreq: changefreq,
		Priority:   strconv.FormatFloat(priority, 'f', 1, 32),
	})
}
