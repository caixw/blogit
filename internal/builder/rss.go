// SPDX-License-Identifier: MIT

package builder

import (
	"html"
	"time"

	"github.com/caixw/blogit/internal/data"
)

const (
	rssVersion    = "2.0"
	rssDateFormat = time.RFC822
)

type rss struct {
	XMLName struct{} `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *channel `xml:"channel"`
}

type channel struct {
	Title         string  `xml:"title"`
	Link          string  `xml:"link"`
	Description   string  `xml:"description"`
	PubDate       string  `xml:"pubDate,omitempty"`
	LastBuildDate string  `xml:"lastBuildDate,omitempty"`
	Items         []*item `xml:"item,omitempty"`
}

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
}

func (b *Builder) buildRSS(d *data.Data) error {
	if d.RSS == nil {
		return nil
	}

	r := &rss{
		Version: rssVersion,
		Channel: &channel{
			Title:         d.RSS.Title,
			Link:          d.URL,
			Description:   d.Subtitle,
			PubDate:       d.Uptime.Format(rssDateFormat),
			LastBuildDate: d.Modified.Format(rssDateFormat),
			Items:         make([]*item, 0, len(d.RSS.Posts)),
		},
	}

	for _, p := range d.RSS.Posts {
		r.Channel.Items = append(r.Channel.Items, &item{
			Title:       p.Title,
			Link:        p.Permalink,
			Description: html.EscapeString(p.Summary),
			PubDate:     p.Created.Format(rssDateFormat),
		})
	}

	return b.appendXMLFile(d, d.RSS.Path, d.RSS.XSLPermalink, r)
}
