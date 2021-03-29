// SPDX-License-Identifier: MIT

package cmd

import (
	"io"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
	"github.com/caixw/blogit/filesystem"
)

var (
	previewSrc  string
	previewBase string
	previewCert string
	previewKey  string
)

// initPreview 注册 preview 子命令
func initPreview(opt *cmdopt.CmdOpt) {
	fs := opt.New("preview", "以预览的方式运行 HTTP 服务\n", preview)
	fs.StringVar(&previewSrc, "src", "./", "指定源码目录")
	fs.StringVar(&previewBase, "base", "http://localhost:8080", "服务基地址")
	fs.StringVar(&previewCert, "cert", "", "https 模式下需要提供的 cert")
	fs.StringVar(&previewKey, "key", "", "https 模式下需要提供的 key")
}

func preview(w io.Writer) error {
	o := &blogit.Options{
		Src:     previewSrc,
		Dest:    filesystem.Memory(),
		BaseURL: previewBase,
		Cert:    previewCert,
		Key:     previewKey,
		Erro:    erro.asLogger(),
		Info:    info.asLogger(),
		Succ:    succ.asLogger(),
	}
	s, err := blogit.Watch(o)
	if err != nil {
		return err
	}
	return s.Serve()
}
