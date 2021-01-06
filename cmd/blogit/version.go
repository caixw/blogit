// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"io"
	"runtime"

	"github.com/caixw/blogit"
	"github.com/issue9/cmdopt"
)

func initVersion(opt *cmdopt.CmdOpt) {
	opt.New("version", "显示版本号", func(w io.Writer) error {
		_, err := fmt.Fprintf(w, "blogit %s build with %s", blogit.Version, runtime.Version())
		return err
	})
}
