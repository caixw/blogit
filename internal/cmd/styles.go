// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma/styles"
	"github.com/issue9/cmdopt"
)

func initStyles(opt *cmdopt.CmdOpt) {
	opt.New("styles", "显示可用的代码高亮名称\n", func(w io.Writer) error {
		names := styles.Names()
		_, err := fmt.Println(names)
		return err
	})
}
