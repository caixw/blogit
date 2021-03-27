// SPDX-License-Identifier: MIT

// 静态博客网站生成工具
//
// 可通过 blogit 查看具体的子命令。
package main

import "github.com/caixw/blogit/internal/cmd"

func main() {
	if err := cmd.Exec(); err != nil {
		panic(err)
	}
}
