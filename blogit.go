// SPDX-License-Identifier: MIT

// Package blogit 静态博客生成工具
//
//
// 本地化
//
// 在 internal/locale 提供了本地化的翻译内容，可以通过 github.com/issue9/localeutil
// 加载相应的翻译内容。
//
// 返回的错误信息，也提供了本地化支持，只要判断该错误对象是否实现了 localeutil.LocaleStringer
// 接口即可，如果实现了，调用 LocaleString() 方法会输出本地的错误信息。
//  b := catalog.NewBuilder()
//  localeutil.LoadMessageFromFile(b, "internal/locale/*.yaml", yaml.Unmarshal)
//  p := message.NewPrinter(language.Chinese, message.Catalog(b))
//
//  err := Build(...)
//  if ls, ok := err.(localeutil.LocaleStringer); ok {
//      println(ls.LocaleString(p)) // 输出本地化内容
//  } else {
//      println(err.Error())
//  }
package blogit

import (
	"io/fs"
	"log"

	"github.com/caixw/blogit/v2/builder"
	"github.com/caixw/blogit/v2/internal/vars"
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
// src 表示源码目录，直接读该文件系统根目录下的内容；
// dest 表示写入的文件系统，默认提供了 DirFS 和 MemoryFS；
// info 输出编译的进度信息，可以为空；
func Build(src fs.FS, dest WritableFS, info *log.Logger) error {
	return NewBuilder(src, dest).Rebuild(info, "")
}

// NewBuilder 声明 Builder
//
// src 表示源码目录，直接读该文件系统根目录下的内容；
// dest 表示写入的文件系统，默认提供了 DirFS 和 MemoryFS；
func NewBuilder(src fs.FS, dest WritableFS) *Builder { return builder.New(src, dest) }

// DirFS 以普通目录结构作为保存对象的文件系统
func DirFS(dir string) WritableFS { return builder.DirFS(dir) }

// MemoryFS 以内在作为保存实体的文件系统
func MemoryFS() WritableFS { return builder.MemoryFS() }
