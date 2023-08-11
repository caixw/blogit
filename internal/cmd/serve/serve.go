// SPDX-License-Identifier: MIT

// Package serve 提供 serve 子命令
package serve

import (
	"flag"
	"io"
	"net/http"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/cmd/console"
)

var opt *options

// Init 注册 serve 子命令
//
// 与 preview 的区别在于，preview 会显示草稿，且可以修改 baseURL，而 serve 不行。
func Init(o *cmdopt.CmdOpt, succ, info, erro *console.Logger, p *message.Printer) {
	o.New("serve", localeutil.StringPhrase("serve title").LocaleString(p), localeutil.StringPhrase("serve usage").LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		opt = &options{p: p}
		fs.StringVar(&opt.source, "src", "./", localeutil.StringPhrase("serve src").LocaleString(p))
		fs.StringVar(&opt.dest, "dest", "", localeutil.StringPhrase("serve dest").LocaleString(p))
		fs.StringVar(&opt.addr, "addr", ":8080", localeutil.StringPhrase("serve port").LocaleString(p))
		fs.StringVar(&opt.path, "path", "/", localeutil.StringPhrase("serve path").LocaleString(p))
		fs.StringVar(&opt.cert, "cert", "", localeutil.StringPhrase("serve http cert").LocaleString(p))
		fs.StringVar(&opt.key, "key", "", localeutil.StringPhrase("serve http key").LocaleString(p))
		fs.StringVar(&opt.hookMethod, "hook.method", http.MethodPost, localeutil.StringPhrase("serve web hook method").LocaleString(p))
		fs.StringVar(&opt.hookURL, "hook.url", "", localeutil.StringPhrase("serve web hook url").LocaleString(p))
		fs.StringVar(&opt.hookAuth, "hook.auth", "", localeutil.StringPhrase("serve web hook auth").LocaleString(p))

		return func(w io.Writer) error {
			if err := opt.serve(succ, info, erro); err != nil {
				if ls, ok := err.(localeutil.LocaleStringer); ok {
					erro.Println(ls.LocaleString(p))
				} else {
					erro.Println(err)
				}
			}
			return nil
		}
	})

}
