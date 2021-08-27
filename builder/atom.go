// SPDX-License-Identifier: MIT

package builder

import (
	"time"

	"github.com/caixw/blogit/v2/internal/data"
)

const (
	atomDateFormat = time.RFC3339
	atomNamespace  = "http://www.w3.org/2005/Atom"
)

type atom struct {
	XMLName  struct{}     `xml:"feed"`
	XMLNS    string       `xml:"xmlns,attr"`
	Title    atomContent  `xml:"title"`
	Subtitle atomContent  `xml:"subtitle"`
	ID       string       `xml:"id"`
	Updated  string       `xml:"updated"`
	Links    []*atomLink  `xml:"link,omitempty"`
	Entries  []*atomEntry `xml:"entry,omitempty"`
}

type atomEntry struct {
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

func (b *Builder) buildAtom(d *data.Data) error {
	if d.Atom == nil {
		return nil
	}

	a := &atom{
		XMLNS:    atomNamespace,
		Title:    atomContent{Content: d.Atom.Title},
		Subtitle: atomContent{Content: d.Subtitle},
		ID:       d.URL,
		Updated:  d.Modified.Format(atomDateFormat),
		Links: []*atomLink{
			{Href: d.URL},
			{Href: d.Atom.Permalink, Rel: "self"},
		},
		Entries: make([]*atomEntry, 0, len(d.Atom.Posts)),
	}

	for _, p := range d.Atom.Posts {
		a.Entries = append(a.Entries, &atomEntry{
			Title:   atomContent{Content: p.Title},
			ID:      p.Permalink,
			Updated: p.Modified.Format(atomDateFormat),
			Links: []*atomLink{
				{Href: p.Permalink, Type: "application/xml"},
			},
			Summary: &atomContent{Type: "text/html", Content: p.Summary},
		})
	}

	return b.appendXMLFile(d.Atom.Path, d.Atom.XSLPermalink, a)
}
