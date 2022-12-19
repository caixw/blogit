// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/caixw/blogit/v2/internal/vars"
	"github.com/issue9/assert/v3"
)

func TestPrintDrafts(t *testing.T) {
	a := assert.New(t, false)

	draftsSrc = "../testdata"
	w := bytes.Buffer{}
	a.NotError(printDrafts(&w))
	a.True(strings.Contains(w.String(), vars.DraftTitleAround+"draft"+vars.DraftTitleAround))
}
