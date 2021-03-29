// SPDX-License-Identifier: MIT

package filesystem

import (
	"errors"
	"io/fs"
)

// Exists 判断文件或是目录是否存在
func Exists(fsys fs.FS, path string) bool {
	_, err := fsys.Open(path)
	return err == nil || errors.Is(err, fs.ErrExist)
}
