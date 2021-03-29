// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"os"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
	"github.com/caixw/blogit/filesystem"
)

var (
	serveSrc  string
	serveAddr string
	servePath string
)

// initServe 注册 serve 子命令
func initServe(opt *cmdopt.CmdOpt) {
	fs := opt.New("serve", "以 HTTP 服务运行\n", serve)
	fs.StringVar(&serveSrc, "src", "./", "指定源码目录")
	fs.StringVar(&serveAddr, "addr", ":8080", "服务端口")
	fs.StringVar(&servePath, "path", "/", "根路径")
}

func serve(w io.Writer) error {
	o := &blogit.Options{
		Src:  os.DirFS(serveSrc),
		Dest: filesystem.Memory(),
		Addr: serveAddr,
		Path: servePath,
		Erro: erro.asLogger(),
		Info: info.asLogger(),
		Succ: succ.asLogger(),
	}
	s, err := blogit.Serve(o)
	if err != nil {
		return err
	}
	return s.Serve()
}
