// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/cmd/console"
)

func TestPrintVersion(t *testing.T) {
	a := assert.New(t, false)

	p, err := console.NewPrinter()
	a.NotError(err).NotNil(p)

	pv := printVersion(p)
	w := bytes.Buffer{}
	a.NotError(pv(&w))
	a.True(strings.Contains(w.String(), blogit.Version(false)))
}
