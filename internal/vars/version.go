// SPDX-License-Identifier: MIT

package vars

const mainVersion = "1.3.0"

var (
	buildDate  string
	commitHash string
	version    = mainVersion
)

func init() {
	if len(buildDate) > 0 {
		version += "+" + buildDate
	}
}

// Version 获取完整的版本号
func Version() string {
	return version
}

// CommitHash 获取最后一条代码提交记录的 hash 值
func CommitHash() string {
	return commitHash
}
