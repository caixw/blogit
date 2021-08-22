// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/caixw/blogit"
)

func TestPrintVersion(t *testing.T) {
	a := assert.New(t)

	w := bytes.Buffer{}
	pv := printVersion(message.NewPrinter(language.Chinese))
	a.NotError(pv(&w))
	a.True(strings.Contains(w.String(), blogit.Version(false)))
}
