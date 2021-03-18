// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/internal/data"
)

const (
	atomVersion    = "1.0"
	atomDateFormat = time.RFC3339
	atomNamespace  = "http://www.w3.org/2005/Atom"
)

type atom struct {
	XMLName  struct{}    `xml:"feed"`
	XMLNS    string      `xml:"xmlns,attr"`
	Title    atomContent `xml:"title"`
	Subtitle atomContent `xml:"subtitle"`
	ID       string      `xml:"id"`
	Updated  string      `xml:"updated"`
	Links    []*atomLink `xml:"link,omitempty"`
	Entries  []*entry    `xml:"entry,omitempty"`
}

type entry struct {
	Title   atomContent  `xml:"title"`
	ID      string       `xml:"id"`
	Updated string       `xml:"updated"`
	Links   []*atomLink  `xml:"link,omitempty"`
	Summary *atomContent `xml:"summary,omitempty"`
}

type atomLink struct {
	Href  string `xml:"href,attr"`
	Type  string `xml:"type,attr,omitempty"`
	Rel   string `xml:"rel,attr,omitempty"`
	Title string `xml:"title,attr,omitempty"`
}

type atomContent struct {
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

func (b *builder) buildAtom(path string, d *data.Data) error {
	if d.Atom == nil {
		return nil
	}

	size := d.RSS.Size
	if len(d.Index.Posts) < size {
		size = len(d.Index.Posts)
	}

	a := &atom{
		XMLNS:    atomNamespace,
		Title:    atomContent{Content: d.Atom.Title},
		Subtitle: atomContent{Content: d.Subtitle},
		ID:       d.URL,
		Updated:  d.Modified.Format(atomDateFormat),
		Links: []*atomLink{
			{Href: d.URL},
			{Href: d.BuildURL(path), Rel: "self"},
		},
		Entries: make([]*entry, 0, size),
	}

	for i := 0; i < size; i++ {
		p := d.Index.Posts[i]
		permalink := d.BuildURL(p.Path)
		a.Entries = append(a.Entries, &entry{
			Title:   atomContent{Content: p.Title},
			ID:      permalink,
			Updated: p.Modified.Format(atomDateFormat),
			Links: []*atomLink{
				{Href: permalink, Type: "application/xml"},
			},
			Summary: &atomContent{Type: "text/html", Content: p.Summary},
		})
	}

	return b.appendXMLFile(d, path, d.Theme.Atom, a)
}
