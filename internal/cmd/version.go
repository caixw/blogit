// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/vars"
)

const (
	versionTitle     = localeutil.StringPhrase("version title")
	versionUsage     = localeutil.StringPhrase("version usage")
	fullVersionUsage = localeutil.StringPhrase("show full version")
)

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("version", versionTitle.LocaleString(p), versionUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var versionFull bool
		fs.BoolVar(&versionFull, "full", false, fullVersionUsage.LocaleString(p))

		return func(w io.Writer) error {
			v := blogit.Version(versionFull)
			_, err := p.Fprintln(w, vars.Name, v, "\n", runtime.Version())
			return err
		}
	})
}
