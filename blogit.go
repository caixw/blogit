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

// Version 返回版本号
//
// full 表示是否返回完整版本号，包含了编译日期，提交的 hash 等额外的值。
func Version(full bool) string {
	if full {
		return vars.FullVersion()
	}
	return vars.Version()
}

// Build 编译并输出内容
//
// dir 表示源码目录，直接读该文件系统根目录下的内容；
// dest 表示输出的目录；
// info 输出编译的进度信息，如果为空，会采用 log.Default()；
func Build(src fs.FS, dest WritableFS, info *log.Logger) error {
	return NewBuilder(dest, info, nil).Rebuild(src, "")
}

func NewBuilder(dest WritableFS, info, erro *log.Logger) *Builder {
	return builder.New(dest, info, erro)
}

func DirFS(dir string) WritableFS {
	return builder.DirFS(dir)
}

func MemoryFS() WritableFS {
	return builder.MemoryFS()
}
