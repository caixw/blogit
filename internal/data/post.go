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
	Path     string
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
			return nil, err
		}
		ps = append(ps, post)
	}

	prevNext(ps)

	return ps, nil
}

func buildPost(conf *loader.Config, theme *loader.Theme, p *loader.Post) (*Post, error) {
	if p.Authors == nil {
		p.Authors = conf.Authors
	}

	if p.License == nil {
		p.License = conf.License
	}

	if p.Language == "" {
		p.Language = conf.Language
	}

	var od *Outdated
	switch p.Outdated {
	case "":
	case loader.OutdatedCreated:
		if conf.Outdated == nil {
			return nil, &loader.FieldError{File: p.Slug, Field: "outdated", Message: "仅允许自定义内容"}
		}
		od = &Outdated{
			Outdated: p.Created.Add(conf.Outdated.Outdated),
			Content:  fmt.Sprintf(conf.Outdated.Created, p.Created.Format(conf.ShortDateFormat)),
		}
	case loader.OutdatedModified:
		if conf.Outdated == nil {
			return nil, &loader.FieldError{File: p.Slug, Field: "outdated", Message: "仅允许自定义内容"}
		}
		od = &Outdated{
			Outdated: p.Modified.Add(conf.Outdated.Outdated),
			Content:  fmt.Sprintf(conf.Outdated.Modified, p.Modified.Format(conf.ShortDateFormat)),
		}
	default:
		od = &Outdated{
			Outdated: time.Time{},
			Content:  p.Outdated,
		}
	}

	if sliceutil.Count(theme.Templates, func(i int) bool { return theme.Templates[i] == p.Template }) == 1 {
		return nil, &loader.FieldError{Message: "不存在", Field: "template", File: p.Slug}
	}

	pp := &Post{
		Slug:     p.Slug,
		Path:     buildPath(p.Slug),
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
		Template: p.Template,
	}

	return pp, nil
}

func prevNext(posts []*Post) {
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
