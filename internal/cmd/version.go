// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit"
)

// initVersion 注册 version 子命令
func initVersion(opt *cmdopt.CmdOpt) {
	opt.New("version", "显示版本号\n", printVersion)
}

func printVersion(w io.Writer) error {
	_, err := fmt.Fprintf(w, "blogit %s\nbuild with %s\n", blogit.FullVersion(), runtime.Version())
	return err
}
