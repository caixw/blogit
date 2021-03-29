// SPDX-License-Identifier: MIT

// Package serve 提供 serve 子命令
package serve

import (
	"io"
	"log"

	"github.com/issue9/cmdopt"
)

var opt = &options{}

// Init 注册 serve 子命令
func Init(o *cmdopt.CmdOpt, info, erro *log.Logger) {
	fs := o.New("serve", "以 HTTP 服务运行\n", func(w io.Writer) error {
		return opt.serve()
	})

	fs.StringVar(&opt.source, "src", "./", "指定源码目录")
	fs.StringVar(&opt.dest, "dest", "", "指定输出目录，如果为空表示采用内存保存。")
	fs.StringVar(&opt.addr, "addr", ":8080", "服务端口")
	fs.StringVar(&opt.path, "path", "/", "项目的访问路径")
	fs.StringVar(&opt.cert, "cert", "", "http 证书")
	fs.StringVar(&opt.key, "key", "", "http 证书")
}
