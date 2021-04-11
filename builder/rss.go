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
	XMLName struct{}    `xml:"rss"`
	Version string      `xml:"version,attr"`
	Channel *rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string     `xml:"title"`
	Link          string     `xml:"link"`
	Description   string     `xml:"description"`
	PubDate       string     `xml:"pubDate,omitempty"`
	LastBuildDate string     `xml:"lastBuildDate,omitempty"`
	Items         []*rssItem `xml:"item,omitempty"`
}

type rssItem struct {
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
		Channel: &rssChannel{
			Title:         d.RSS.Title,
			Link:          d.URL,
			Description:   d.Subtitle,
			PubDate:       d.Uptime.Format(rssDateFormat),
			LastBuildDate: d.Modified.Format(rssDateFormat),
			Items:         make([]*rssItem, 0, len(d.RSS.Posts)),
		},
	}

	for _, p := range d.RSS.Posts {
		r.Channel.Items = append(r.Channel.Items, &rssItem{
			Title:       p.Title,
			Link:        p.Permalink,
			Description: html.EscapeString(p.Summary),
			PubDate:     p.Created.Format(rssDateFormat),
		})
	}

	return b.appendXMLFile(d.RSS.Path, d.RSS.XSLPermalink, r)
}
