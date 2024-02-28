// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package builder

import (
	"log"
	"net/http"
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/assert/v4/rest"

	"github.com/caixw/blogit/v2/internal/testdata"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestIsIgnore(t *testing.T) {
	a := assert.New(t, false)

	a.True(isIgnore("abc/def/a.md"))
	a.True(isIgnore("abc/def/.git"))
	a.True(isIgnore("themes/d/layout/header.html"))
	a.True(isIgnore("themes/layout/layout/header.html")) // 第一个 layout 为主题名称
	a.False(isIgnore("themes/d/theme.yaml"))
	a.False(isIgnore("themes/layout/header.html")) // 不符合 themes/xx/layout 的格式
}

func TestBuilder_Handler(t *testing.T) {
	a := assert.New(t, false)

	// MemoryFS

	b := &Builder{
		Src:     testdata.Source,
		Dest:    MemoryFS(),
		BaseURL: "http://localhost:8080",
	}
	srv := rest.NewServer(a, b.Handler(nil), nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do(nil).Status(http.StatusNotFound)

	a.NotError(b.Rebuild()).
		NotZero(b.Builded())
	srv.Get("/robots.txt").Do(nil).Status(http.StatusOK)
	srv.Get("/posts/p1" + vars.Ext).Do(nil).Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do(nil).Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do(nil).Status(http.StatusOK)

	// index.html
	srv.Get("/").Do(nil).Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do(nil).Status(http.StatusOK)
	srv.Get("/themes/").Do(nil).Status(http.StatusNotFound) // 目录下没有 index.html

	// DirFS

	destDir, err := testdata.Temp()
	a.NotError(err)
	b = &Builder{
		Src:     testdata.Source,
		Dest:    DirFS(destDir),
		Info:    log.Default(),
		BaseURL: "http://localhost:8080",
	}
	srv = rest.NewServer(a, b.Handler(nil), nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do(nil).Status(http.StatusNotFound)

	a.NotError(b.Rebuild())
	srv.Get("/robots.txt").Do(nil).Status(http.StatusOK)
	srv.Get("/posts/p1" + vars.Ext).Do(nil).Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do(nil).Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do(nil).Status(http.StatusOK)

	// index.html
	srv.Get("/").Do(nil).Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do(nil).Status(http.StatusOK)
	srv.Get("/themes/").Do(nil).Status(http.StatusNotFound) // 目录下没有 index.html
}
