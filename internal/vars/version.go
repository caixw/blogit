// SPDX-License-Identifier: MIT

package vars

var (
	version     = "2.3.2" // 版本号，由 goreleaser 负责在编译时更新到最新的 git tag
	metadata    string
	fullVersion = version
)

func init() {
	if metadata != "" {
		fullVersion += "+" + metadata
	}
}

func FullVersion() string { return fullVersion }

func Version() string { return version }
