// SPDX-License-Identifier: MIT

// Package console 输出到控制台的日志
package console

import (
	"io"
	"log"

	"github.com/issue9/term/v2/colors"
)

// Logger 输出到控制台的日志
type Logger struct {
	colors.Colorize
	Prefix string
	Out    io.Writer
}

func (msg *Logger) Printf(format string, v ...interface{}) {
	msg.Fprint(msg.Out, msg.Prefix)
	colors.Fprintf(msg.Out, colors.Normal, colors.Default, colors.Default, format, v...)
}

func (msg *Logger) Println(v ...interface{}) {
	msg.Fprint(msg.Out, msg.Prefix)
	colors.Fprintln(msg.Out, colors.Normal, colors.Default, colors.Default, v...)
}

func (msg *Logger) Write(bs []byte) (int, error) {
	msg.Fprint(msg.Out, msg.Prefix)
	colors.Fprint(msg.Out, colors.Normal, colors.Default, colors.Default, string(bs))
	return 0, nil
}

func (msg *Logger) AsLogger() *log.Logger {
	return log.New(msg, "", log.Ldate)
}
