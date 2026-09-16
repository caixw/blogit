// SPDX-FileCopyrightText: 2020-2026 caixw
//
// SPDX-License-Identifier: MIT

package cmd

import (
	"strings"
	"testing"

	"github.com/issue9/assert/v5"
)

func TestPrintStyles(t *testing.T) {
	a := assert.New(t, false)

	opt, buf, p := newCMD(a)
	initStyles(opt, p)
	a.NotError(opt.Exec([]string{"styles"}))
	a.True(strings.Contains(buf.String(), "solarized-dark256"))
}
