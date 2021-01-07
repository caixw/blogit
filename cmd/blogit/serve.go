// SPDX-License-Identifier: MIT

package main

import (
	"io"
	"log"

	"github.com/caixw/blogit"
	"github.com/issue9/cmdopt"
)

var (
	serveSrc   string
	serveAddr  string
	servePath  string
	serveWatch bool
)

func initServe(opt *cmdopt.CmdOpt) {
	fs := opt.New("serve", "以 HTTP 服务运行", serve)
	fs.StringVar(&serveSrc, "src", "./", "指定源码目录")
	fs.StringVar(&serveAddr, "addr", ":8080", "服务端口")
	fs.StringVar(&servePath, "path", "/", "根路径")
	fs.BoolVar(&serveWatch, "watch", false, "监视变化")
}

func serve(w io.Writer) error {
	if serveWatch {
		return blogit.Watch(serveSrc, serveAddr, servePath, log.New(w, "", log.Ldate))
	}
	return blogit.Serve(serveSrc, serveAddr, servePath)
}
