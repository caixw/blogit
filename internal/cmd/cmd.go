// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
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

const (
	cmdUsage  = localeutil.StringPhrase("cmd usage")
	helpUsage = localeutil.StringPhrase("help usage")
)

// Exec 执行命令行
func Exec(args []string) error {
	systag, _ := localeutil.DetectUserLanguageTag() // 即使出错，依然会返回 language.Tag
	p, err := console.NewPrinter(systag)
	if err != nil {
		return err
	}

	opt := cmdopt.New(os.Stdout, flag.ExitOnError, cmdUsage.LocaleString(p), nil, func(name string) string {
		return localeutil.Phrase("sub command not found %s", name).LocaleString(p)
	})

	initDrafts(opt, p)
	initBuild(opt, p)
	initVersion(opt, p)
	initStyles(opt, p)
	serve.Init(opt, succ, info, erro, p)
	preview.Init(opt, succ, info, erro, p)
	create.InitInit(opt, erro, p)
	create.InitPost(opt, succ, erro, p)
	cmdopt.Help(opt, "help", helpUsage.LocaleString(p), helpUsage.LocaleString(p))

	return opt.Exec(args)
}
