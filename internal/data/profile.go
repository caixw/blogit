// SPDX-License-Identifier: MIT

package data

import "github.com/caixw/blogit/internal/loader"

// Profile github.com 下与账号同名仓库的 README.md 文件管理
type Profile struct {
	Path   string
	Title  string
	Footer string
	Posts  []*Post
}

func newProfile(conf *loader.Config, posts []*Post) *Profile {
	p := conf.Profile

	size := p.Size
	if s := len(posts); s < size {
		size = s
	}

	profile := &Profile{
		Path:   "README.md",
		Title:  p.Title,
		Footer: p.Footer,
		Posts:  make([]*Post, 0, size),
	}

	for i := 0; i < size; i++ {
		profile.Posts = append(profile.Posts, posts[i])
	}

	return profile
}
