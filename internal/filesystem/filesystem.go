// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package filesystem 提供文件系统的相关函数
package filesystem

import (
	"errors"
	"io/fs"
)

// Exists 判断文件或是目录是否存在
func Exists(f fs.FS, path string) bool {
	_, err := fs.Stat(f, path)
	return err == nil || errors.Is(err, fs.ErrExist)
}
