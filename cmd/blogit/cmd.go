// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/issue9/cmdopt"
)

func main() {
	opt := &cmdopt.CmdOpt{
		Output:        os.Stdout,
		ErrorHandling: flag.ExitOnError,
		Header:        "静态博客工具",
		Footer:        "源码以 MIT 许可发布于 https://github.com/caixw/blogit",
		CommandsTitle: "子命令",
		OptionsTitle:  "参数",
		NotFound: func(name string) string {
			return fmt.Sprintf("未找到子命令 %s", name)
		},
	}

	initBuild(opt)
	initServe(opt)
	initVersion(opt)
	opt.Help("help", "显示当前内容")

	if err := opt.Exec(os.Args[1:]); err != nil {
		panic(err)
	}
}
