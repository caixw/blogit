// SPDX-License-Identifier: MIT

package loader

import (
	"io/fs"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/v2/internal/testdata"
	"github.com/caixw/blogit/v2/locale"
)

func newPrinter() (*message.Printer, error) {
	systag, _ := localeutil.DetectUserLanguageTag() // 即使出错，依然会返回 language.Tag

	b := catalog.NewBuilder()
	if err := localeutil.LoadMessageFromFSGlob(b, locale.Locales(), "*.yaml", yaml.Unmarshal); err != nil {
		return nil, err
	}

	return message.NewPrinter(systag, message.Catalog(b)), nil
}

func TestTag_sanitize(t *testing.T) {
	a := assert.New(t, false)
	tag := &Tag{}
	a.Error(tag.sanitize(nil))

	tag.Slug = "s1"
	a.Error(tag.sanitize(nil))

	tag.Title = "t1"
	a.Error(tag.sanitize(nil))

	tag.Content = "c1"
	a.NotError(tag.sanitize(nil))
	a.NotError(tag.sanitize(&Tags{Tags: []*Tag{tag}}))

	p, err := newPrinter()
	a.NotError(err).NotNil(p)
	e := tag.sanitize(&Tags{Tags: []*Tag{tag, tag}})
	a.Contains(e.LocaleString(p), p.Sprintf("duplicate value"))
}

func TestLoadTags(t *testing.T) {
	a := assert.New(t, false)

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
