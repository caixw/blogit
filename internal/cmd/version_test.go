// SPDX-License-Identifier: MIT

package cmd

import (
	"strings"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2"
)

func TestPrintVersion(t *testing.T) {
	a := assert.New(t, false)

	opt, buf, p := newCMD(a)
	initVersion(opt, p)
	a.NotError(opt.Exec([]string{"version"}))
	a.True(strings.Contains(buf.String(), blogit.Version(false)))
}
