// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestProfile_sanitize(t *testing.T) {
	a := assert.New(t)

	p := &Profile{}
	err := p.sanitize()
	a.Equal(err.Field, "title")

	p.Title = "# title"
	err = p.sanitize()
	a.Equal(err.Field, "size")

	p.Size = 5
	p.Footer = "#### \tfooter"
	err = p.sanitize()
	a.NotError(err)

	a.Equal(p.Title, "### title").Equal(p.Footer, "##### footer")

	// p.Alternate !=""

	p = &Profile{
		Alternate: "path/to/file",
		Title:     "title",
		Size:      5,
		Footer:    "footer",
	}

	err = p.sanitize()
	a.Equal(err.Field, "title")

	p.Title = ""
	err = p.sanitize()
	a.Equal(err.Field, "size")

	p.Size = 0
	err = p.sanitize()
	a.Equal(err.Field, "footer")

	p.Footer = ""
	err = p.sanitize()
	a.NotError(err)
}
