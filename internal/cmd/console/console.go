// SPDX-License-Identifier: MIT

// Package console 输出到控制台的日志
package console

import (
	"log"

	"github.com/issue9/term/v3/colors"
)

// Logger 输出到控制台的日志
type Logger struct {
	Colorize *colors.Colorize
	Prefix   string
	Color    colors.Color
}

func (l *Logger) printPrefix() *colors.Colorize {
	return l.Colorize.Color(colors.Normal, l.Color, colors.Default).Print(l.Prefix).Reset()
}

func (l *Logger) Printf(format string, v ...interface{}) { l.printPrefix().Printf(format, v...) }

func (l *Logger) Println(v ...interface{}) { l.printPrefix().Println(v...) }

func (l *Logger) Write(bs []byte) (int, error) {
	l.printPrefix().Print(string(bs))
	return len(bs), nil
}

func (l *Logger) AsLogger() *log.Logger {
	return log.New(l, "", log.Ldate)
}
