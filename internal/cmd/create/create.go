// SPDX-License-Identifier: MIT

// Package create 创建文章的相关子命令
package create

import (
	"os"

	"github.com/caixw/blogit/filesystem"
)

func getWD() (filesystem.WritableFS, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return filesystem.Dir(dir), nil
}
