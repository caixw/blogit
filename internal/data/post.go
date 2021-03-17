// SPDX-License-Identifier: MIT

package data

import (
	"sort"
	"time"

	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Post 文章详情
type Post struct {
	Permalink string
	Slug      string
	Path      string
	Title     string
	Created   time.Time
	Modified  time.Time
	Tags      []*Tag
	tags      []string
	Language  string
	Authors   []*loader.Author
	License   *loader.License
	Summary   string
	Content   string
	Image     string
	Prev      *Post
	Next      *Post
	Template  string
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

	if p.Modified.IsZero() {
		p.Modified = p.Created
	}

	if p.Template == "" {
		p.Template = vars.DefaultTemplate
	}

	if sliceutil.Count(theme.Templates, func(i int) bool { return theme.Templates[i] == p.Template }) == 0 {
		return nil, &loader.FieldError{Message: "模板不存在于 theme.yaml", Field: "template", File: p.Slug + ".md", Value: p.Template}
	}

	path := buildPath(p.Slug)
	pp := &Post{
		Permalink: buildURL(conf.URL, path),
		Slug:      p.Slug,
		Path:      path,
		Title:     p.Title,
		Created:   p.Created,
		Modified:  p.Modified,
		tags:      p.Tags,
		Language:  p.Language,
		Authors:   p.Authors,
		License:   p.License,
		Summary:   p.Summary,
		Content:   p.Content,
		Image:     p.Image,
		Template:  p.Template,
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
