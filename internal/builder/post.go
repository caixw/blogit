// SPDX-License-Identifier: MIT

package builder

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
	Authors  []*Author
	License  *License
	Content  string
	Prev     *Post
	Next     *Post
}

// Outdated 定义文章过时显示的信息
type Outdated struct {
	Outdated time.Time
	Content  string
}

func buildPosts(conf *loader.Config, posts []*loader.Post) ([]*Post, error) {
	if err := checkRawPosts(posts); err != nil {
		return nil, err
	}

	ps := make([]*Post, 0, len(posts))
	for _, p := range posts {
		post, err := buildPost(conf, p)
		if err != nil {
			return nil, err
		}
		ps = append(ps, post)
	}

	prevNext(ps)

	return ps, nil
}

func buildPost(conf *loader.Config, p *loader.Post) (*Post, error) {
	if p.Authors == nil {
		p.Authors = conf.Authors
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
			Content:  p.Content,
		}
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
		Content:  p.Content,
		Slug:     p.Slug,
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

// 检测原始的文章内容，并对其进行排序。
func checkRawPosts(posts []*loader.Post) *loader.FieldError {
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

	for _, p := range posts {
		cnt := sliceutil.Count(posts, func(i int) bool {
			return p.Slug == posts[i].Slug && p.Slug != posts[i].Slug
		})
		if cnt > 1 {
			return &loader.FieldError{Message: "存在重复的值", Field: "slug"}
		}
	}

	return nil
}
