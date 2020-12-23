// SPDX-License-Identifier: MIT

package loader

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)
	conf, err := loadConfig("./testdata/conf.yaml")
	a.NotError(err).NotNil(conf)

	a.Equal(conf.Author.Name, "caixw")
	a.Equal(conf.License.Rel, "license")
	a.Equal(conf.Language, "cmn-Hans")

	conf, err = loadConfig("./testdata/not-exists.yaml")
	a.ErrorIs(err, os.ErrNotExist).Nil(conf)

	conf, err = loadConfig("./testdata/failed_conf.yaml")
	a.ErrorType(err, &FieldError{}, err).Nil(conf)
}

func TestConfig_BuildURL(t *testing.T) {
	a := assert.New(t)
	conf, err := loadConfig("./testdata/conf.yaml")
	a.NotError(err).NotNil(conf)

	a.Equal(conf.BuildURL("/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(conf.BuildURL("p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(conf.BuildURL(""), "https://example.com/")
	a.Equal(conf.BuildURL("/"), "https://example.com/")
}

func TestArchive_sanitize(t *testing.T) {
	a := assert.New(t)

	archive := &Archive{}
	a.Error(archive.sanitize())

	archive = &Archive{Format: "2001 年"}
	a.NotError(archive.sanitize())
	a.Equal(archive.Order, ArchiveOrderDesc).Equal(archive.Type, ArchiveTypeYear)

	archive = &Archive{Format: "2001 年", Type: "not-exists"}
	a.Equal(archive.sanitize().Field, "type")

	archive = &Archive{Format: "2001 年", Type: ArchiveTypeMonth, Order: "not-exists"}
	a.Equal(archive.sanitize().Field, "order")
}

func TestRSS_sanitize(t *testing.T) {
	a := assert.New(t)

	rss := &RSS{}
	conf := &Config{
		Title: "title",
		RSS:   rss,
	}
	a.Error(rss.sanitize(conf))

	// Size 错误
	rss.Size = 0
	a.Error(rss.sanitize(conf))
	rss.Size = -1
	a.Error(rss.sanitize(conf))

	rss.Size = 10
	a.NotError(rss.sanitize(conf))
}

func TestSitemap_sanitize(t *testing.T) {
	a := assert.New(t)

	s := &Sitemap{}
	a.Error(s.sanitize())

	s.Priority = -1.0
	a.Error(s.sanitize())
	s.Priority = 1.1
	a.Error(s.sanitize())

	s.Priority = .8
	s.PostPriority = 0.9
	s.Changefreq = "never"
	s.PostChangefreq = "never"
	a.NotError(s.sanitize())
}
