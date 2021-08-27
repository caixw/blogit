// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"runtime"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/vars"
)

var versionFull bool

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt, p *message.Printer) {
	fs := opt.New("version", p.Sprintf("version usage"), printVersion(p))
	fs.BoolVar(&versionFull, "full", false, p.Sprintf("show full version"))
}

func printVersion(p *message.Printer) func(io.Writer) error {
	return func(w io.Writer) error {
		v := blogit.Version(versionFull)
		_, err := p.Fprintf(w, "version content", vars.Name, v, runtime.Version())
		return err
	}
}
