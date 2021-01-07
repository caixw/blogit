// SPDX-License-Identifier: MIT

package data

import (
	"path"
	"time"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Tag 单个标签的内容
type Tag struct {
	Slug     string
	Path     string
	Title    string
	Content  string // 对该标签的详细描述
	Posts    []*Post
	Created  time.Time
	Modified time.Time
}

func buildTags(tags []*loader.Tag) ([]*Tag, error) {
	ts := make([]*Tag, 0, len(tags))
	for _, t := range tags {
		ts = append(ts, &Tag{
			Slug:    t.Slug,
			Path:    buildPath(path.Join(vars.TagsDir, t.Slug)),
			Title:   t.Title,
			Content: t.Content,
		})
	}

	return ts, nil
}

func checkTags(tags []*Tag, posts []*Post) (created, modified time.Time, err error) {
	for _, p := range posts {
		if created.Before(p.Created) {
			created = p.Created
		}
		if modified.Before(p.Modified) {
			modified = p.Modified
		}

		for _, tag := range p.tags {
			t := findTagByName(tags, tag)
			if t == nil {
				return time.Time{}, time.Time{}, &loader.FieldError{File: p.Slug, Message: "不存在", Field: "tags." + tag}
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
	}

	if modified.IsZero() {
		modified = created
	}
	return created, modified, nil
}

func findTagByName(tags []*Tag, slug string) *Tag {
	for _, t := range tags {
		if t.Slug == slug {
			return t
		}
	}
	return nil
}
