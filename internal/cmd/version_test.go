// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert"
)

func TestPrintVersion(t *testing.T) {
	a := assert.New(t)

	w := bytes.Buffer{}
	a.NotError(printVersion(&w))
	a.True(strings.Contains(w.String(), "solarized-dark256"))
}
