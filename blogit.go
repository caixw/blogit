// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
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

// Version 版本号
func Version() string {
	return vars.Version()
}

// Build 编译并输出内容
//
// dir 表示源码目录；
// dest 表示输出的目录；
func Build(src fs.FS, dest WritableFS) error {
	return NewBuilder(dest, nil).Rebuild(src, "")
}

func NewBuilder(dest WritableFS, erro *log.Logger) *Builder {
	return builder.New(dest, erro)
}
