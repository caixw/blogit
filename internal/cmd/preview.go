// SPDX-License-Identifier: MIT

package cmd

import (
	"io"

	"github.com/caixw/blogit"
	"github.com/issue9/cmdopt"
)

var (
	previewSrc  string
	previewBase string
)

// initPreview 注册 preview 子命令
func initPreview(opt *cmdopt.CmdOpt) {
	fs := opt.New("preview", "以预览的方式运行 HTTP 服务\n", preview)
	fs.StringVar(&previewSrc, "src", "./", "指定源码目录")
	fs.StringVar(&previewBase, "base", "http://localhost:8080", "服务基地址")
}

func preview(w io.Writer) error {
	return blogit.Watch(previewSrc, previewBase, info.asLogger(), erro.asLogger(), succ.asLogger())
}
