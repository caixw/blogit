// SPDX-License-Identifier: MIT

package data

import (
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/loader"
)

func TestLoad(t *testing.T) {
	a := assert.New(t)

	data, err := Load("../testdata")
	a.NotError(err).NotNil(data)

	a.Equal(data.Icon.Type, "image/png").Equal(data.Icon.Sizes, "256x256")
	a.Equal(5, len(data.Menus))
	a.Equal(3, len(data.Posts)).
		Equal(data.Posts[1].Prev, data.Posts[0]).
		Equal(data.Posts[1].Next, data.Posts[2])
	a.NotNil(data.Posts[1].Authors)
	a.NotNil(data.Posts[1].License)
	a.NotNil(data.License)
	a.NotEmpty(data.Authors)

	a.True(data.Builded.After(time.Time{}))
}

func TestData_BuildURL(t *testing.T) {
	a := assert.New(t)

	data := &Data{URL: "https://example.com/"} // 传入的 URL 必定是以 / 结尾的
	a.Equal(data.BuildURL("/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(data.BuildURL("p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(data.BuildURL(""), "https://example.com/")
	a.Equal(data.BuildURL("/"), "https://example.com/")
}

func TestData_buildThemeURL(t *testing.T) {
	a := assert.New(t)

	data := &Data{URL: "https://example.com/", Theme: &loader.Theme{ID: "def"}} // 传入的 URL 必定是以 / 结尾的
	a.Equal(data.BuildThemeURL("/p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(data.BuildThemeURL("p1/p2.md"), "https://example.com/themes/def/p1/p2.md")
	a.Equal(data.BuildThemeURL(""), "https://example.com/themes/def")
	a.Equal(data.BuildThemeURL("/"), "https://example.com/themes/def")
}
