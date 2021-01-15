// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/issue9/cmdopt"
)

// 定义了命令行的输出通道
var (
	infoWriter = os.Stdout
	erroWriter = os.Stderr
	succWriter = os.Stdout
)

// Exec 执行命令行
func Exec() error {
	opt := &cmdopt.CmdOpt{
		Output:        infoWriter,
		ErrorHandling: flag.ExitOnError,
		Header:        "静态博客工具\n",
		Footer:        "源码以 MIT 许可发布于 https://github.com/caixw/blogit\n",
		CommandsTitle: "子命令",
		OptionsTitle:  "参数",
		NotFound: func(name string) string {
			return fmt.Sprintf("未找到子命令 %s\n", name)
		},
	}

	initBuild(opt)
	initServe(opt)
	initInit(opt)
	initPost(opt)
	initVersion(opt)
	opt.Help("help", "显示当前内容")

	return opt.Exec(os.Args[1:])
}
