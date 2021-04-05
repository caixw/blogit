// SPDX-License-Identifier: MIT

package vars

// Version 主版本号
const Version = "1.5.0" // NOTE: 应该保持与 tag 同值

var (
	metadata    string
	fullVersion = Version
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

// FullVersion 获取完整的版本号
func FullVersion() string {
	return fullVersion
}
