// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package data

import "github.com/caixw/blogit/v2/internal/loader"

// RSS 整理后的 RSS 和 Atom 数据
type RSS struct {
	Title        string
	Permalink    string
	XSLPermalink string
	Path         string
	Posts        []*Post
}

func newRSS(conf *loader.Config, r *loader.RSS, path, xsl string, posts []*Post) *RSS {
	size := r.Size
	if l := len(posts); l < size {
		size = l
	}

	rss := &RSS{
		Title:     r.Title,
		Permalink: BuildURL(conf.URL, path),
		Path:      path,
		Posts:     make([]*Post, 0, size),
	}

	if xsl != "" {
		rss.XSLPermalink = buildThemeURL(conf.URL, conf.Theme, xsl)
	}

	for i := 0; i < size; i++ {
		rss.Posts = append(rss.Posts, posts[i])
	}

	return rss
}
