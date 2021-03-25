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

	p.Title = "### title"
	err = p.sanitize()
	a.Equal(err.Field, "size")

	p.Size = 5
	err = p.sanitize()
	a.NotError(err)

	// p.Alternate !=""

	p.Alternate = "path/to/file"
	err = p.sanitize()
	a.Equal(err.Field, "title")

	p.Title = ""
	err = p.sanitize()
	a.Equal(err.Field, "size")

	p.Size = 0
	err = p.sanitize()
	a.NotError(err)
}
