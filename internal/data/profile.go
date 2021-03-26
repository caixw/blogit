// SPDX-License-Identifier: MIT

package data

import "github.com/caixw/blogit/internal/loader"

// Profile github.com 下与账号同名仓库的 README.md 文件管理
type Profile struct {
	*loader.Profile
	Path string
}

func newProfile(conf *loader.Config) *Profile {
	return &Profile{
		Path:    buildPath("README.md"),
		Profile: conf.Profile,
	}
}
