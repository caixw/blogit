// SPDX-License-Identifier: MIT

package data

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
)

// Index 索引页内容
type Index struct {
	Title       string
	Permalink   string
	Keywords    string
	Description string
	Posts       []*Post
	Index       int // 当前页的索引
	Path        string
	Next        *Index
	Prev        *Index
}

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
	License   *loader.Link
	Keywords  string
	Summary   string
	Content   string
	Image     string
	Prev      *Post
	Next      *Post
	Template  string
	JSONLD    string
	TOC       []loader.Header
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

	postsPrevNext(ps)

	return ps, nil
}

func buildIndexes(conf *loader.Config, posts []*Post) []*Index {
	size := int(math.Ceil(float64(len(posts)) / float64(conf.Index.Size)))
	indexes := make([]*Index, 0, size)
	hasVerbs := strings.Contains(conf.Index.Title, "%d")

	for page := 0; page < size; page++ {
		start := page * conf.Index.Size
		offset := start + conf.Index.Size
		if offset > len(posts) {
			offset = len(posts)
		}

		index := &Index{
			Keywords:    conf.Keywords,
			Description: conf.Description,
			Posts:       posts[start:offset],
			Index:       page + 1,
		}
		if hasVerbs {
			index.Title = fmt.Sprintf(conf.Index.Title, index.Index)
		} else {
			index.Title = conf.Index.Title
		}

		if page == 0 {
			index.Path = vars.IndexFilename
			index.Permalink = conf.URL
		} else {
			index.Path = fmt.Sprintf(vars.IndexFilenameFormat, index.Index)
			index.Permalink = BuildURL(conf.URL, index.Path)
		}

		indexes = append(indexes, index)
	}

	for i, index := range indexes {
		if i > 0 {
			index.Prev = indexes[i-1]
		}
		if i < len(indexes)-1 {
			index.Next = indexes[i+1]
		}
	}

	return indexes
}

func buildPost(conf *loader.Config, theme *loader.Theme, p *loader.Post) (*Post, error) {
	if p.Authors == nil {
		p.Authors = []*loader.Author{conf.Author}
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

	if len(p.TOC) <= conf.TOC {
		p.TOC = nil
	}

	if sliceutil.Count(theme.Templates, func(i int) bool { return theme.Templates[i] == p.Template }) == 0 {
		return nil, &loader.FieldError{
			Message: localeutil.Phrase("template not found in", vars.ThemeYAML),
			Field:   "template",
			File:    p.Slug + vars.MarkdownExt,
			Value:   p.Template,
		}
	}

	// NOTE: p.JSONLD 用到以上的一些变量，比如 p.License 等，所以需要放在最后初始化。
	if p.JSONLD == "" {
		ld, err := buildPostLD(p)
		if err != nil {
			return nil, err
		}
		p.JSONLD = ld
	}

	path := p.Slug + vars.Ext
	return &Post{
		Permalink: BuildURL(conf.URL, path),
		Slug:      p.Slug,
		Path:      path,
		Title:     p.Title,
		Created:   p.Created,
		Modified:  p.Modified,
		tags:      p.Tags,
		Language:  p.Language,
		Authors:   p.Authors,
		License:   p.License,
		Keywords:  p.Keywords,
		Summary:   p.Summary,
		Content:   p.Content,
		Image:     p.Image,
		Template:  p.Template,
		JSONLD:    p.JSONLD,
		TOC:       p.TOC,
	}, nil
}

func postsPrevNext(posts []*Post) {
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

func sortPostsByCreated(posts []*Post) []*Post {
	sorted := make([]*Post, len(posts))
	copy(sorted, posts)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Created.After(sorted[j].Created)
	})

	return sorted
}
