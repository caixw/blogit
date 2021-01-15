// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"log"

	"github.com/issue9/term/v2/colors"
)

var (
	erro = &consoleLogger{
		prefix:   "[ERRO] ",
		Colorize: colors.New(colors.Normal, colors.Red, colors.Default),
		out:      erroWriter,
	}

	info = &consoleLogger{
		prefix:   "[INFO] ",
		Colorize: colors.New(colors.Normal, colors.Default, colors.Default),
		out:      infoWriter,
	}

	succ = &consoleLogger{
		prefix:   "[SUCC] ",
		Colorize: colors.New(colors.Normal, colors.Green, colors.Default),
		out:      succWriter,
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
