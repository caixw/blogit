// SPDX-License-Identifier: MIT

// Package utils 提供通用函数
package utils

import "os"

// FileExists 判断文件或是目录是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
