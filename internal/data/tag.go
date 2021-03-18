// SPDX-License-Identifier: MIT

package data

import (
	"path"
	"time"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
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
	Content   string // 对该标签的详细描述
	Posts     []*Post
	Created   time.Time
	Modified  time.Time
}

func buildTags(conf *loader.Config, tags *loader.Tags) (*Tags, error) {
	if tags.Keywords == "" {
		tags.Keywords = conf.Keywords
	}
	if tags.Description == "" {
		tags.Description = conf.Description
	}

	ts := &Tags{
		Title:       tags.Title,
		Permalink:   buildURL(conf.URL, vars.TagsFilename),
		Keywords:    tags.Keywords,
		Description: tags.Description,
	}
	for _, t := range tags.Tags {
		p := buildPath(path.Join(vars.TagsDir, t.Slug))
		ts.Tags = append(ts.Tags, &Tag{
			Permalink: buildURL(conf.URL, p),
			Slug:      t.Slug,
			Path:      p,
			Title:     t.Title,
			Content:   t.Content,
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
