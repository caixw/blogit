// SPDX-License-Identifier: MIT

package vars

var (
	version     = "1.6.0" // 版本号，由 goreleaser 负责在编译时更新到最新的 git tag
	metadata    string
	fullVersion = version
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

// FullVersion 获取完整的版本号
func FullVersion() string { return fullVersion }

// Version 返回版本信息
func Version() string { return version }
