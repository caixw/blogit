// SPDX-License-Identifier: MIT

package data

import (
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/v2/internal/testdata"
)

func TestLoad(t *testing.T) {
	a := assert.New(t)

	data, err := Load(testdata.Source, false, "")
	a.NotError(err).NotNil(data)

	a.Equal(data.Icon.Type, "image/png").Equal(data.Icon.Sizes, "256x256")
	a.Equal(3, len(data.Posts)).
		Equal(data.Posts[1].Prev, data.Posts[0]).
		Equal(data.Posts[1].Next, data.Posts[2])
	a.NotNil(data.Posts[1].Authors)
	a.NotNil(data.Posts[1].License)
	a.NotNil(data.License)
	a.NotNil(data.Author)
	a.Equal(2, len(data.Indexes)) // 3 篇文章，每页 2 篇，可分为 2 个索引页
	a.Equal(data.URL, "https://example.com")

	a.True(data.Builded.After(time.Time{}))

	data, err = Load(testdata.Source, true, "https://example.com/v2")
	a.NotError(err).NotNil(data)
	a.Equal(data.URL, "https://example.com/v2")
	a.Equal(4, len(data.Posts))
}

func TestBuildURL(t *testing.T) {
	a := assert.New(t)

	base := "https://example.com/"
	a.Equal(BuildURL(base, "/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(BuildURL(base, "p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(BuildURL(base, ""), "https://example.com/")
	a.Equal(BuildURL(base, "/"), "https://example.com/")

	base = "https://example.com"
	a.Equal(BuildURL(base, "/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(BuildURL(base, "p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(BuildURL(base, ""), "https://example.com/")
	a.Equal(BuildURL(base, "/"), "https://example.com/")

	base = ""
	a.Equal(BuildURL(base, "/p1/p2.md"), "/p1/p2.md")
	a.Equal(BuildURL(base, "p1/p2.md"), "/p1/p2.md")
	a.Equal(BuildURL(base, ""), "/")
	a.Equal(BuildURL(base, "/"), "/")
}

func TestBuildThemeURL(t *testing.T) {
	a := assert.New(t)

	url := "https://example.com/"
	id := "def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com"
	id = "def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com/"
	id = "/def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com"
	id = "/def"
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes/def")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes/def")

	url = "https://example.com/"
	id = ""
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes")

	url = "https://example.com"
	id = ""
	a.Equal(buildThemeURL(url, id, "/p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, "p1/p2.md"), "https://example.com/themes/p1/p2.md")
	a.Equal(buildThemeURL(url, id, ""), "https://example.com/themes")
	a.Equal(buildThemeURL(url, id, "/"), "https://example.com/themes")
}
