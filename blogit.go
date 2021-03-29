// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
package blogit

import (
	"io/fs"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/vars"
)

// Version 版本号
func Version() string {
	return vars.Version()
}

// Build 编译并输出内容
//
// dir 表示源码目录；
// dest 表示输出的目录；
func Build(src fs.FS, dest filesystem.WritableFS) error {
	return builder.Build(src, dest)
}
