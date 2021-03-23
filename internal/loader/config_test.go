// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)
	conf, err := LoadConfig("../../testdata/src/conf.yaml")
	a.NotError(err).NotNil(conf)

	a.Equal(conf.Author.Name, "caixw")
	a.Equal(conf.Language, "cmn-Hans")

	conf, err = LoadConfig("../../testdata/src/not-exists.yaml")
	a.ErrorIs(err, os.ErrNotExist).Nil(conf)

	conf, err = LoadConfig("../../testdata/src/failed_conf.yaml")
	a.ErrorType(err, &FieldError{}, err).Nil(conf)
}

func TestRSS_sanitize(t *testing.T) {
	a := assert.New(t)

	rss := &RSS{}
	conf := &Config{
		Title: "title",
		RSS:   rss,
	}
	a.Error(rss.sanitize())

	// Size 错误
	rss.Size = 0
	a.Error(rss.sanitize())
	rss.Size = -1
	a.Error(rss.sanitize())

	rss.Size = 10
	a.NotError(rss.sanitize())
	a.Equal(rss.Title, conf.Title)
}
