// SPDX-FileCopyrightText: 2020-2026 caixw
//
// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"

	"github.com/issue9/assert/v5"
	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/language"

	"github.com/caixw/blogit/v2/internal/cmd/console"
)

func newCMD(a *assert.Assertion) (*cmdopt.CmdOpt, *bytes.Buffer, *localeutil.Printer) {
	buf := &bytes.Buffer{}
	opt := cmdopt.New(&cmdopt.Options{
		Output:        buf,
		ErrorHandling: flag.ContinueOnError,
	})
	a.NotNil(opt)

	p, err := console.NewPrinter(language.SimplifiedChinese)
	a.NotError(err).NotNil(p)

	return opt, buf, p
}
