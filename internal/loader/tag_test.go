// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/testdata"
)

func TestTag_sanitize(t *testing.T) {
	a := assert.New(t)
	tag := &Tag{}
	a.Error(tag.sanitize(nil))

	tag.Slug = "s1"
	a.Error(tag.sanitize(nil))

	tag.Title = "t1"
	a.Error(tag.sanitize(nil))

	tag.Content = "c1"
	a.NotError(tag.sanitize(nil))
	a.NotError(tag.sanitize(&Tags{Tags: []*Tag{tag}}))
	a.ErrorString(tag.sanitize(&Tags{Tags: []*Tag{tag, tag}}), "重复的值")
}

func TestLoadTags(t *testing.T) {
	a := assert.New(t)

	tags, err := LoadTags(testdata.Source, "tags.yaml")
	a.NotError(err).NotNil(tags).Equal(tags.Title, "标签").Equal(tags.OrderType, TagOrderTypeSize)
	a.Equal(4, len(tags.Tags))
	a.Equal(tags.Tags[0].Slug, "default").
		Equal(tags.Tags[1].Slug, "api").
		Equal(tags.Tags[2].Slug, "firefox").
		Equal(tags.Tags[3].Slug, "git")

	tags, err = LoadTags(testdata.Source, "not-exists.yaml")
	a.ErrorIs(err, fs.ErrNotExist).Empty(tags)
}
