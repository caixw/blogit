// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestTag_sanitize(t *testing.T) {
	a := assert.New(t)
	data := &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("conf.yaml"))

	tag := &Tag{}
	a.Error(tag.sanitize(nil, data.Config))

	tag.Slug = "s1"
	a.Error(tag.sanitize(nil, data.Config))

	tag.Title = "t1"
	a.Error(tag.sanitize(nil, data.Config))

	tag.Content = "c1"
	a.NotError(tag.sanitize(nil, data.Config))
	a.NotError(tag.sanitize([]*Tag{tag}, data.Config))
	a.ErrorString(tag.sanitize([]*Tag{tag, tag}, data.Config), "重复的值")
}

func TestLoadTags(t *testing.T) {
	a := assert.New(t)

	data := &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("conf.yaml"))
	a.NotError(data.loadTags("./tags.yaml")).NotNil(data.Tags)
	a.Equal(4, len(data.Tags))
	a.Equal(data.Tags[0].Slug, "default").
		Equal(data.Tags[1].Slug, "api").
		Equal(data.Tags[2].Slug, "firefox").
		Equal(data.Tags[3].Slug, "git")

	data = &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("conf.yaml"))
	a.ErrorIs(data.loadTags("./not-exists.yaml"), os.ErrNotExist).Empty(data.Tags)

	data = &Data{Dir: "./testdata"}
	a.NotError(data.loadConfig("conf.yaml"))
	err := data.loadTags("./failed_tags.yaml")
	a.Error(err).Empty(data.Tags)
	ferr, ok := err.(*FieldError)
	a.True(ok).Equal(ferr.File, "testdata/failed_tags.yaml")
}
