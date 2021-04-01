// SPDX-License-Identifier: MIT

// Package blogit 静态博客生成工具
package blogit

import (
	"io/fs"
	"log"

	"github.com/caixw/blogit/builder"
	"github.com/caixw/blogit/internal/vars"
)

type (
	Builder    = builder.Builder
	WritableFS = builder.WritableFS
)

// Version 当前的版本号
const Version = vars.Version

// FullVersion 返回完整的版本号
///
// 完整版本号包含了编译日期，提交的 hash 等额外的值。
func FullVersion() string { return vars.FullVersion() }

// Build 编译并输出内容
//
// dir 表示源码目录，直接读该文件系统根目录下的内容；
// dest 表示输出的目录；
func Build(src fs.FS, dest WritableFS) error {
	return NewBuilder(dest, nil).Rebuild(src, "")
}

func NewBuilder(dest WritableFS, erro *log.Logger) *Builder {
	return builder.New(dest, erro)
}
