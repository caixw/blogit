// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/testdata"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	conf, err := LoadConfig(testdata.Source, "conf.yaml")
	a.NotError(err).NotNil(conf)

	a.Equal(conf.Author.Name, "author1")
	a.Equal(conf.Language, "cmn-Hans")

	conf, err = LoadConfig(testdata.Source, "not-exists.yaml")
	a.ErrorIs(err, os.ErrNotExist).Nil(conf)
}

func TestRSS_sanitize(t *testing.T) {
	a := assert.New(t)

	rss := &RSS{}
	err := rss.sanitize()
	a.Equal(err.Field, "title")

	// Size 错误
	rss.Title = "title"
	rss.Size = 0
	err = rss.sanitize()
	a.Equal(err.Field, "size")
	rss.Size = -1
	err = rss.sanitize()
	a.Equal(err.Field, "size")

	rss.Size = 10
	a.NotError(rss.sanitize())
}
