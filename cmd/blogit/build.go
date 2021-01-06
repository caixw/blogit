// SPDX-License-Identifier: MIT

package main

import (
	"io"

	"github.com/caixw/blogit"
	"github.com/issue9/cmdopt"
)

var buildDir string

func initBuild(opt *cmdopt.CmdOpt) {
	fs := opt.New("build", "编译内容", build)
	fs.StringVar(&buildDir, "dir", "./", "指定源码目录")
}

func build(w io.Writer) error {
	return blogit.Build(buildDir)
}
