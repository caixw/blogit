// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// 静态博客网站生成工具
//
// 可通过 blogit help 查看具体的子命令。
package main

import (
	"os"

	"github.com/caixw/blogit/v2/internal/cmd"
)

func main() {
	if err := cmd.Exec(os.Args[1:]); err != nil {
		panic(err)
	}
}
