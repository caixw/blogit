// SPDX-License-Identifier: MIT

package testdata

import (
	"embed"
	"os"
)

//go:embed posts themes conf.yaml tags.yaml
var Source embed.FS

// Temp 创建一个临时的文件夹
func Temp() (string, error) {
	return os.MkdirTemp(os.TempDir(), "blogit")
}
