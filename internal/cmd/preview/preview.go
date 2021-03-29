// SPDX-License-Identifier: MIT

// Package preview 提供 preview 子命令
package preview

import (
	"io"
	"log"

	"github.com/issue9/cmdopt"
)

var opt = &options{}

// Init 注册 preview 子命令
func Init(o *cmdopt.CmdOpt, succ, info, erro *log.Logger) {
	opt.succ = succ
	opt.info = info
	opt.erro = erro

	fs := o.New("preview", "以预览的方式运行 HTTP 服务\n", func(w io.Writer) error {
		return opt.watch()
	})

	fs.StringVar(&opt.source, "src", "./", "指定源码目录")
	fs.StringVar(&opt.dest, "dest", "", "指定保存了对象，为空表示保存在内存。")
	fs.StringVar(&opt.url, "url", "http://localhost:8080", "服务基地址")
	fs.StringVar(&opt.cert, "cert", "", "https 模式下需要提供的 cert")
	fs.StringVar(&opt.key, "key", "", "https 模式下需要提供的 key")
}
