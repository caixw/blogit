// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package preview 提供 preview 子命令
package preview

import (
	"flag"
	"io"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/cmd/console"
)

var opt *options

// Init 注册 preview 子命令
func Init(o *cmdopt.CmdOpt, succ, info, erro *console.Logger, p *message.Printer) {
	opt = &options{p: p}
	o.New("preview", localeutil.StringPhrase("preview title").LocaleString(p), localeutil.StringPhrase("preview usage").LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		fs.StringVar(&opt.source, "src", "./", localeutil.StringPhrase("preview src").LocaleString(p))
		fs.StringVar(&opt.dest, "dest", "", localeutil.StringPhrase("preview dest").LocaleString(p))
		fs.StringVar(&opt.url, "url", "http://localhost:8080", localeutil.StringPhrase("preview base url").LocaleString(p))
		fs.StringVar(&opt.cert, "cert", "", localeutil.StringPhrase("preview http cert").LocaleString(p))
		fs.StringVar(&opt.key, "key", "", localeutil.StringPhrase("preview http key").LocaleString(p))
		return func(w io.Writer) error {
			if err := opt.watch(succ, info, erro); err != nil {
				erro.Println(err)
			}
			return nil
		}
	})

}
