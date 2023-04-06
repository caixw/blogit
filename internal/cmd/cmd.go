// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v3/colors"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/cmd/create"
	"github.com/caixw/blogit/v2/internal/cmd/preview"
	"github.com/caixw/blogit/v2/internal/cmd/serve"
)

var (
	erro = &console.Logger{
		Prefix:   "[ERRO] ",
		Colorize: colors.New(os.Stderr),
		Color:    colors.Red,
	}

	info = &console.Logger{
		Prefix:   "[INFO] ",
		Colorize: colors.New(os.Stdout),
		Color:    colors.Yellow,
	}

	succ = &console.Logger{
		Prefix:   "[SUCC] ",
		Colorize: colors.New(os.Stdout),
		Color:    colors.Green,
	}
)

// Exec 执行命令行
func Exec(args []string) error {
	p, err := console.NewPrinter()
	if err != nil {
		return err
	}

	usage := "cmd usage"
	opt := cmdopt.New(os.Stdout, flag.ExitOnError, usage, nil, func(name string) string {
		return p.Sprintf("sub command not found", name)
	})

	initDrafts(opt, p)
	initBuild(opt, p)
	initVersion(opt, p)
	initStyles(opt, p)
	serve.Init(opt, succ, info, erro, p)
	preview.Init(opt, succ, info, erro, p)
	create.InitInit(opt, erro, p)
	create.InitPost(opt, succ, erro, p)
	opt.New("help", p.Sprintf("help usage"), p.Sprintf("help usage"), cmdopt.Help(opt))

	return opt.Exec(args)
}
