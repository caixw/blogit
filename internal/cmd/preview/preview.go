// SPDX-License-Identifier: MIT

// Package preview 提供 preview 子命令
package preview

import (
	"io"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/internal/cmd/console"
)

var opt = &options{}

// Init 注册 preview 子命令
func Init(o *cmdopt.CmdOpt, succ, info, erro *console.Logger, p *message.Printer) {
	fs := o.New("preview", p.Sprintf("preview usage"), func(w io.Writer) error {
		if err := opt.watch(succ, info, erro); err != nil {
			erro.Println(err)
		}
		return nil
	})

	fs.StringVar(&opt.source, "src", "./", p.Sprintf("preview src"))
	fs.StringVar(&opt.dest, "dest", "", p.Sprintf("preview dest"))
	fs.StringVar(&opt.url, "url", "http://localhost:8080", p.Sprintf("preview base url"))
	fs.StringVar(&opt.cert, "cert", "", p.Sprintf("preview http cert"))
	fs.StringVar(&opt.key, "key", "", p.Sprintf("preview http key"))
}
