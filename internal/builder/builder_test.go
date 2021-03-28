// SPDX-License-Identifier: MIT

package builder

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/blogit/internal/vars"
)

func TestIsIgnore(t *testing.T) {
	a := assert.New(t)

	a.True(isIgnore("abc/def/a.md"))
	a.True(isIgnore("abc/def/.git"))
	a.True(isIgnore("abc/themes/d/layout/header.html"))
	a.False(isIgnore("abc/themes/layout/header.html")) // 不符合 themes/xx/layout 的格式
}

func TestBuild(t *testing.T) {
	a := assert.New(t)

	a.NotError(os.RemoveAll("../../testdata/dest/index.xml"))

	err := Build("../../testdata/src", "../../testdata/dest")
	a.NotError(err)
	a.FileExists("../../testdata/dest/index" + vars.Ext).
		FileExists("../../testdata/dest/tags" + vars.Ext).
		FileExists("../../testdata/dest/tags/default" + vars.Ext).
		FileExists("../../testdata/dest/posts/p1" + vars.Ext)
}

func TestBuilder_appendFile(t *testing.T) {
	a := assert.New(t)
	now := time.Now()

	b := &Builder{files: make([]*file, 0, 10)}
	a.Panic(func() {
		b.appendFile("", now, []byte("<html><head></head></html>"))
	})

	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("abc.html", now, []byte("#h1\n\n##h2"))
	a.NotNil(b.files[0].data).Equal(b.files[0].path, "abc.html")

	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("/abc.html", now, []byte("#h1\n\n##h2"))
	a.NotNil(b.files[0].data).Equal(b.files[0].path, "abc.html")

}

func TestBuilder_ServeHTTP(t *testing.T) {
	a := assert.New(t)

	b := &Builder{}
	srv := rest.NewServer(t, b, nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do().Status(http.StatusNotFound)

	a.NotError(b.Build("../../testdata/src", "http://localhost:8080"))
	srv.Get("/robots.txt").Do().Status(http.StatusOK)
	srv.Get("/posts/p1" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do().Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do().Status(http.StatusOK)

	// index.html
	srv.Get("/").Do().Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do().Status(http.StatusOK)
	srv.Get("/themes/").Do().Status(http.StatusNotFound) // 目录下没有 index.html
}
