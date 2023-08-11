// SPDX-License-Identifier: MIT

package data

import (
	"path"
	"sort"
	"strings"
	"time"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
)

// Tags 标签列表及相关设置项
type Tags struct {
	Title       string
	Permalink   string
	Keywords    string
	Description string
	Tags        []*Tag
}

// Tag 单个标签的内容
type Tag struct {
	Permalink string
	Slug      string
	Path      string
	Title     string
	Keywords  string
	Content   string // 对该标签的详细描述
	Posts     []*Post
	Prev      *Tag
	Next      *Tag
	Created   time.Time
	Modified  time.Time
}

func buildTags(conf *loader.Config, tags *loader.Tags, ps []*Post) (*Tags, error) {
	if tags.Keywords == "" {
		tags.Keywords = conf.Keywords
	}
	if tags.Description == "" {
		tags.Description = conf.Description
	}

	ts := &Tags{
		Title:       tags.Title,
		Permalink:   BuildURL(conf.URL, vars.TagsFilename),
		Keywords:    tags.Keywords,
		Description: tags.Description,
		Tags:        make([]*Tag, 0, len(tags.Tags)),
	}

	keys := make([]string, 0, len(tags.Tags))
	for _, t := range tags.Tags {
		keys = append(keys, t.Slug)

		key := t.Slug
		if t.Slug != t.Title {
			key += "," + t.Title
			keys = append(keys, t.Title)
		}

		p := path.Join(vars.TagsDir, t.Slug+vars.Ext)
		ts.Tags = append(ts.Tags, &Tag{
			Permalink: BuildURL(conf.URL, p),
			Slug:      t.Slug,
			Path:      p,
			Title:     t.Title,
			Content:   t.Content,
			Keywords:  key,
		})
	}

	if ts.Keywords == "" {
		ts.Keywords = strings.Join(keys, ",")
	}

	if err := ts.relationTagsPosts(ps); err != nil {
		return nil, err
	}
	ts.clearTags() // 清除无文章关联的标签
	sortTags(ts.Tags, tags.OrderType, tags.Order)
	tagsPrevNext(ts.Tags)

	return ts, nil
}

func sortTags(tags []*Tag, typ, order string) {
	if typ == loader.TagOrderTypeSize {
		sort.SliceStable(tags, func(i, j int) bool {
			return len(tags[i].Posts) > len(tags[j].Posts)
		})
	}

	if order == loader.OrderDesc {
		sliceutil.Reverse(tags)
	}
}

// 关联 tags 和 posts 的信息
func (ts *Tags) relationTagsPosts(posts []*Post) error {
	for _, p := range posts {
		for _, tag := range p.tags {
			t := findTagByName(ts.Tags, tag)
			if t == nil {
				return &loader.FieldError{File: p.Slug, Message: localeutil.Phrase("not found"), Field: "tags." + tag}
			}
			t.Posts = append(t.Posts, p)
			p.Tags = append(p.Tags, t)

			if t.Created.Before(p.Created) {
				t.Created = p.Created
			}

			if t.Modified.Before(p.Modified) {
				t.Modified = p.Modified
			}
		}

		if p.Keywords == "" {
			keywords := make([]string, 0, len(p.Tags)*2)
			for _, t := range p.Tags {
				keywords = append(keywords, t.Slug, t.Title)
			}
			keywords = sliceutil.Unique(keywords, func(i, j string) bool { return i == j })
			p.Keywords = strings.Join(keywords, ",")
		}
	}

	return nil
}

func findTagByName(tags []*Tag, slug string) *Tag {
	for _, t := range tags {
		if t.Slug == slug {
			return t
		}
	}
	return nil
}

func (ts *Tags) clearTags() {
	ts.Tags = sliceutil.Delete(ts.Tags, func(i *Tag, _ int) bool { return len(i.Posts) == 0 })
}

func tagsPrevNext(tags []*Tag) {
	max := len(tags)
	for i := 0; i < max; i++ {
		tag := tags[i]
		if i > 0 {
			tag.Prev = tags[i-1]
		}
		if i < max-1 {
			tag.Next = tags[i+1]
		}
	}
}
