// SPDX-License-Identifier: MIT

package cmd

import (
	"strings"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/blogit/v2/internal/vars"
)

func TestPrintDrafts(t *testing.T) {
	a := assert.New(t, false)

	opt, buf, p := newCMD(a)
	initDrafts(opt, p)
	a.NotError(opt.Exec([]string{"drafts", "-src", "../testdata"}))

	a.True(strings.Contains(buf.String(), vars.DraftTitleAround+"draft"+vars.DraftTitleAround))
}
