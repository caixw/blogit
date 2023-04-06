// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
)

// initBuild 注册 build 子命令
func initBuild(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("build", p.Sprintf("build title"), p.Sprintf("build usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var buildSrc string
		var buildDest string
		fs.StringVar(&buildSrc, "src", "./", p.Sprintf("build src"))
		fs.StringVar(&buildDest, "dest", "./dest", p.Sprintf("build dest"))

		return func(w io.Writer) error {
			start := time.Now()

			info.Println(p.Sprintf("start build"))
			if err := blogit.Build(os.DirFS(buildSrc), blogit.DirFS(buildDest), info.AsLogger()); err != nil {
				if ls, ok := err.(localeutil.LocaleStringer); ok {
					erro.Println(ls.LocaleString(p))
				} else {
					erro.Println(err)
				}
				return nil
			}

			succ.Println(p.Sprintf("build complete", time.Since(start)))
			return nil
		}
	})
}
