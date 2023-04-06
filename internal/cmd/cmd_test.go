// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"

	"github.com/issue9/assert/v3"
	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"

	"github.com/caixw/blogit/v2/internal/cmd/console"
)

func newCMD(a *assert.Assertion) (*cmdopt.CmdOpt, *bytes.Buffer, *localeutil.Printer) {
	buf := &bytes.Buffer{}
	opt := cmdopt.New(buf, flag.ContinueOnError, "", nil, nil)
	a.NotNil(opt)

	p, err := console.NewPrinter()
	a.NotError(err).NotNil(p)

	return opt, buf, p
}
