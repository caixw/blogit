// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma/styles"
	"github.com/issue9/cmdopt"
)

const stylesUsage = `显示可用的代码高亮名称

代码高亮的相关样式表是以文件的方式关联到 HTML 的，
这些文件将以 themes/{id}/{name}.css 格式出现，
其中 {id} 是主题名称，面 {name} 则是当前命令列的值。`

func initStyles(opt *cmdopt.CmdOpt) {
	opt.New("styles", stylesUsage, printStyles)
}

func printStyles(w io.Writer) error {
	_, err := fmt.Fprintln(w, styles.Names())
	return err
}
