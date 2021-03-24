// SPDX-License-Identifier: MIT

package loader

// Profile 用于生成 github.com/profile 中的 README.md 内容
type Profile struct {
	Alternate string `yaml:"alternate"` // 采用此文件的内容代替

	// 当 alternate 为空，以下值才生效
	Modified int    `yaml:"modified"`         // 显示最后修改的 n 条记录
	Created  int    `yaml:"created"`          // 显示最后添加的 n 条记录
	Header   string `yaml:"header,omitempty"` // 页眉
	Footer   string `yaml:"footer,omitempty"` // 页脚
}
