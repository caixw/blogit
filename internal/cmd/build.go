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

const (
	buildTitle     = localeutil.StringPhrase("build title")
	buildUsage     = localeutil.StringPhrase("build usage")
	buildSrcUsage  = localeutil.StringPhrase("build src")
	buildDestUsage = localeutil.StringPhrase("build dest")
)

// initBuild 注册 build 子命令
func initBuild(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("build", buildTitle.LocaleString(p), buildUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var buildSrc string
		var buildDest string
		fs.StringVar(&buildSrc, "src", "./", buildSrcUsage.LocaleString(p))
		fs.StringVar(&buildDest, "dest", "./dest", buildDestUsage.LocaleString(p))

		return func(w io.Writer) error {
			start := time.Now()

			info.Println(localeutil.StringPhrase("start build").LocaleString(p))
			if err := blogit.Build(os.DirFS(buildSrc), blogit.DirFS(buildDest), info.AsLogger()); err != nil {
				if ls, ok := err.(localeutil.LocaleStringer); ok {
					erro.Println(ls.LocaleString(p))
				} else {
					erro.Println(err)
				}
				return nil
			}

			succ.Println(localeutil.StringPhrase("build complete").LocaleString(p), time.Since(start))
			return nil
		}
	})
}
