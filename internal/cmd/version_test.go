// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit"
	"github.com/caixw/blogit/internal/locale"
)

func TestPrintVersion(t *testing.T) {
	a := assert.New(t)

	p, err := locale.NewPrinter()
	a.NotError(err).NotNil(p)

	pv := printVersion(p)
	w := bytes.Buffer{}
	a.NotError(pv(&w))
	a.True(strings.Contains(w.String(), blogit.Version(false)))
}
