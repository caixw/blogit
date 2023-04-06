// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/vars"
)

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("version", p.Sprintf("version title"), p.Sprintf("version usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var versionFull bool
		fs.BoolVar(&versionFull, "full", false, p.Sprintf("show full version"))

		return func(w io.Writer) error {
			v := blogit.Version(versionFull)
			_, err := p.Fprintf(w, "version content", vars.Name, v, runtime.Version())
			return err
		}
	})
}
