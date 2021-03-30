// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v2/colors"

	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/cmd/create"
	"github.com/caixw/blogit/internal/cmd/preview"
	"github.com/caixw/blogit/internal/cmd/serve"
)

var (
	erro = &console.Logger{
		Prefix:   "[ERRO] ",
		Colorize: colors.New(colors.Normal, colors.Red, colors.Default),
		Out:      os.Stderr,
	}

	info = &console.Logger{
		Prefix:   "[INFO] ",
		Colorize: colors.New(colors.Normal, colors.Default, colors.Default),
		Out:      os.Stdout,
	}

	succ = &console.Logger{
		Prefix:   "[SUCC] ",
		Colorize: colors.New(colors.Normal, colors.Green, colors.Default),
		Out:      os.Stdout,
	}
)

// Exec 执行命令行
func Exec(args []string) error {
	opt := &cmdopt.CmdOpt{
		Output:        os.Stdout,
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
	initVersion(opt)
	serve.Init(opt, info, erro)
	preview.Init(opt, succ, info, erro)
	create.InitInit(opt, succ, erro)
	create.InitPost(opt, succ, erro)
	opt.Help("help", "显示当前内容\n")

	return opt.Exec(args)
}
