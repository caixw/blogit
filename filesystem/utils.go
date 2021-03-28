// SPDX-License-Identifier: MIT

package filesystem

import "os"

// Exists 判断文件或是目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
