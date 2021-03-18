// SPDX-License-Identifier: MIT

package builder

import (
	"html"
	"time"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
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

func (b *builder) buildRSS(d *data.Data) error {
	if d.RSS == nil {
		return nil
	}

	size := d.RSS.Size
	if len(d.Index.Posts) < size {
		size = len(d.Index.Posts)
	}

	r := &rss{
		Version: rssVersion,
		Channel: &channel{
			Title:         d.RSS.Title,
			Link:          d.URL,
			Description:   d.Subtitle,
			PubDate:       d.Uptime.Format(rssDateFormat),
			LastBuildDate: d.Modified.Format(rssDateFormat),
			Items:         make([]*item, 0, size),
		},
	}

	for i := 0; i < size; i++ {
		p := d.Index.Posts[i]
		r.Channel.Items = append(r.Channel.Items, &item{
			Title:       p.Title,
			Link:        p.Permalink,
			Description: html.EscapeString(p.Summary),
			PubDate:     p.Created.Format(rssDateFormat),
		})
	}

	return b.appendXMLFile(d, vars.RssXML, d.RSS.XSL, r)
}
