// SPDX-License-Identifier: MIT

package data

import (
	"fmt"
	"sort"
	"time"

	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/internal/loader"
)

// Post 文章详情
type Post struct {
	Slug     string
	Title    string
	Created  time.Time
	Modified time.Time
	Tags     []*Tag
	tags     []string
	Language string
	Outdated *Outdated
	Authors  []*loader.Author
	License  *loader.License
	Summary  string
	Content  string
	Prev     *Post
	Next     *Post
	Template string
}

// Outdated 定义文章过时显示的信息
type Outdated struct {
	Outdated time.Time
	Content  string
}

func buildPosts(conf *loader.Config, theme *loader.Theme, posts []*loader.Post) ([]*Post, error) {
	sortPosts(posts)

	ps := make([]*Post, 0, len(posts))
	for _, p := range posts {
		post, err := buildPost(conf, theme, p)
		if err != nil {
			err.File = p.Slug
			return nil, err
		}
		ps = append(ps, post)
	}

	prevNext(ps)

	return ps, nil
}

func buildPost(conf *loader.Config, theme *loader.Theme, p *loader.Post) (*Post, *loader.FieldError) {
	if p.Authors == nil {
		p.Authors = conf.Authors
	}

	if p.License == nil {
		p.License = conf.License
	}

	var od *Outdated
	switch p.Outdated {
	case loader.OutdatedCreated:
		od = &Outdated{
			Outdated: p.Created.Add(conf.Outdated),
			Content:  fmt.Sprintf("当前文章创建于 %s，可能已经过时！", p.Created.Format(conf.ShortDateFormat)),
		}
	case loader.OutdatedModified:
		od = &Outdated{
			Outdated: p.Modified.Add(conf.Outdated),
			Content:  fmt.Sprintf("当前文章最后次修改于 %s，可能已经过时！", p.Modified.Format(conf.ShortDateFormat)),
		}
	default:
		od = &Outdated{
			Outdated: p.Created,
			Content:  p.Outdated,
		}
	}

	if sliceutil.Count(theme.Templates, func(i int) bool { return theme.Templates[i] == p.Template }) == 1 {
		return nil, &loader.FieldError{Message: "不存在", Field: "template"}
	}

	pp := &Post{
		Title:    p.Title,
		Created:  p.Created,
		Modified: p.Modified,
		tags:     p.Tags,
		Language: p.Language,
		Outdated: od,
		Authors:  p.Authors,
		License:  p.License,
		Summary:  p.Summary,
		Content:  p.Content,
		Slug:     p.Slug,
		Template: p.Template,
	}

	return pp, nil
}

func prevNext(posts []*Post) {
	// 生成 prev 和 next
	max := len(posts)
	for i := 0; i < max; i++ {
		post := posts[i]
		if i > 0 {
			post.Prev = posts[i-1]
		}
		if i < max-1 {
			post.Next = posts[i+1]
		}
	}
}

func sortPosts(posts []*loader.Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		switch {
		case (posts[i].State == loader.StateTop) || (posts[j].State == loader.StateLast):
			return true
		case (posts[i].State == loader.StateLast) || (posts[j].State == loader.StateTop):
			return false
		default:
			return posts[i].Created.After(posts[j].Created)
		}
	})
}
