// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
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
	a.NotError(tag.sanitize([]*Tag{tag}))
	a.ErrorString(tag.sanitize([]*Tag{tag, tag}), "重复的值")
}

func TestLoadTags(t *testing.T) {
	a := assert.New(t)

	tags, err := LoadTags("../testdata/tags.yaml")
	a.NotError(err).NotNil(tags)
	a.Equal(4, len(tags))
	a.Equal(tags[0].Slug, "default").
		Equal(tags[1].Slug, "api").
		Equal(tags[2].Slug, "firefox").
		Equal(tags[3].Slug, "git")

	tags, err = LoadTags("../testdata/not-exists.yaml")
	a.ErrorIs(err, os.ErrNotExist).Empty(tags)

	tags, err = LoadTags("../testdata/failed_tags.yaml")
	a.Error(err).Nil(tags)
	ferr, ok := err.(*FieldError)
	a.True(ok).Equal(ferr.File, "../testdata/failed_tags.yaml")
}
