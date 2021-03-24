// SPDX-License-Identifier: MIT

// Package cmd 提供命令行相关的功能
package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v2/colors"
)

var (
	erro = &consoleLogger{
		prefix:   "[ERRO] ",
		Colorize: colors.New(colors.Normal, colors.Red, colors.Default),
		out:      os.Stderr,
	}

	info = &consoleLogger{
		prefix:   "[INFO] ",
		Colorize: colors.New(colors.Normal, colors.Default, colors.Default),
		out:      os.Stdout,
	}

	succ = &consoleLogger{
		prefix:   "[SUCC] ",
		Colorize: colors.New(colors.Normal, colors.Green, colors.Default),
		out:      os.Stdout,
	}
)

type consoleLogger struct {
	colors.Colorize
	prefix string
	out    io.Writer
}

func (msg *consoleLogger) printf(format string, v ...interface{}) {
	msg.Fprint(msg.out, msg.prefix)
	colors.Fprintf(msg.out, colors.Normal, colors.Default, colors.Default, format, v...)
}

func (msg *consoleLogger) println(v ...interface{}) {
	msg.Fprint(msg.out, msg.prefix)
	colors.Fprintln(msg.out, colors.Normal, colors.Default, colors.Default, v...)
}

func (msg *consoleLogger) Write(bs []byte) (int, error) {
	msg.Fprint(msg.out, msg.prefix)
	colors.Fprint(msg.out, colors.Normal, colors.Default, colors.Default, string(bs))
	return 0, nil
}

func (msg *consoleLogger) asLogger() *log.Logger {
	return log.New(msg, "", log.Ldate)
}

// Exec 执行命令行
func Exec() error {
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
	initServe(opt)
	initPreview(opt)
	initInit(opt)
	initPost(opt)
	initVersion(opt)
	opt.Help("help", "显示当前内容")

	return opt.Exec(os.Args[1:])
}
