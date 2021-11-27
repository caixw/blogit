// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/blogit/v2/internal/testdata"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t, false)

	conf, err := LoadConfig(testdata.Source, "conf.yaml")
	a.NotError(err).NotNil(conf)

	a.Equal(conf.Author.Name, "author1")
	a.Equal(conf.Language, "cmn-Hans")

	conf, err = LoadConfig(testdata.Source, "not-exists.yaml")
	a.ErrorIs(err, fs.ErrNotExist).Nil(conf)
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t, false)

	conf := &Config{}
	err := conf.sanitize()
	a.Error(err).Equal(err.Field, "url")

	conf = &Config{URL: "https://example.com"}
	err = conf.sanitize()
	a.Error(err).Equal(err.Field, "uptime")

	conf = &Config{
		URL:     "https://example.com",
		Uptime:  time.Now(),
		Title:   "title",
		Theme:   "default",
		Author:  &Author{Name: "example", URL: "https://example.com"},
		Index:   &Index{Title: "%d page", Size: 5},
		Archive: &Archive{Title: "archive"},
		License: &Link{Text: "MIT", URL: "https://example.com"},
	}
	err = conf.sanitize()
	a.NotError(err).Equal(conf.Language, "cmn-Hans")

	conf.Atom = &RSS{}
	err = conf.sanitize()
	a.Error(err).Equal(err.Field, "atom.title")
}

func TestRSS_sanitize(t *testing.T) {
	a := assert.New(t, false)

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

func TestIndex_sanitize(t *testing.T) {
	a := assert.New(t, false)

	i := &Index{}
	err := i.sanitize()
	a.Equal(err.Field, "size")

	i.Size = 5
	err = i.sanitize()
	a.Equal(err.Field, "title")

	i.Title = "xx"
	a.NotError(i.sanitize())
}
