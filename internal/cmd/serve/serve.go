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
	o.New("serve", p.Sprintf("serve title"), p.Sprintf("serve usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		opt = &options{p: p}
		fs.StringVar(&opt.source, "src", "./", p.Sprintf("serve src"))
		fs.StringVar(&opt.dest, "dest", "", p.Sprintf("serve dest"))
		fs.StringVar(&opt.addr, "addr", ":8080", p.Sprintf("serve port"))
		fs.StringVar(&opt.path, "path", "/", p.Sprintf("serve path"))
		fs.StringVar(&opt.cert, "cert", "", p.Sprintf("serve http cert"))
		fs.StringVar(&opt.key, "key", "", p.Sprintf("serve http key"))
		fs.StringVar(&opt.hookMethod, "hook.method", http.MethodPost, p.Sprintf("serve web hook method"))
		fs.StringVar(&opt.hookURL, "hook.url", "", p.Sprintf("serve web hook url"))
		fs.StringVar(&opt.hookAuth, "hook.auth", "", p.Sprintf("serve web hook auth"))

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
