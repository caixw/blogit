// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
)

var versionFull bool

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt) {
	fs := opt.New("version", "显示版本号\n", printVersion)
	fs.BoolVar(&versionFull, "full", false, "显示完整的版本号信息")
}

func printVersion(w io.Writer) error {
	v := blogit.Version(versionFull)
	_, err := fmt.Fprintf(w, "blogit %s\nbuild with %s\n", v, runtime.Version())
	return err
}
