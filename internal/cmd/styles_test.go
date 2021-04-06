// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert"
)

func TestPrintStyles(t *testing.T) {
	a := assert.New(t)

	w := bytes.Buffer{}
	a.NotError(printStyles(&w))
	a.True(strings.Contains(w.String(), "solarized-dark256"))
}
