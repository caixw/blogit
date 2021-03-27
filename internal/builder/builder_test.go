// SPDX-License-Identifier: MIT

package builder

import (
	"net/http"
	"os"
	"testing"

	"github.com/caixw/blogit/internal/vars"
	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
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
	a.FileExists("../../testdata/dest/index.html")
}

func TestBuilder_appendFile(t *testing.T) {
	a := assert.New(t)

	b := &Builder{files: make([]*file, 0, 10)}
	a.Panic(func() {
		b.appendFile("", "", []byte("<html><head></head></html>"))
	})

	// 根据后缀名判断
	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("abc.html", "", []byte("#h1\n\n##h2"))
	a.Equal(b.files[0].ct, "text/html").NotNil(b.files[0].data).Equal(b.files[0].path, "abc.html")

	// 自定义
	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("abc.html", "custom", []byte("<html><head></head></html>"))
	a.Equal(b.files[0].ct, "custom").NotNil(b.files[0].data).Equal(b.files[0].path, "abc.html")

	// 根据 data 判断
	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("abc", "", []byte("<html><head></head></html>"))
	a.Equal(b.files[0].ct, "text/html").NotNil(b.files[0].data).Equal(b.files[0].path, "abc")

	b = &Builder{files: make([]*file, 0, 10)}
	b.appendFile("abc", "", []byte("#h1\n\n##h2"))
	a.Equal(b.files[0].ct, "text/plain").NotNil(b.files[0].data).Equal(b.files[0].path, "abc")
}

func TestBuilder_ServeHTTP(t *testing.T) {
	a := assert.New(t)

	b := &Builder{}
	srv := rest.NewServer(t, b, nil)

	// b 未加载任何数据。返回都是 404
	srv.Get("/robots.txt").Do().Status(http.StatusNotFound)

	a.NotError(b.Build("../../testdata/src", "http://localhost:8080"))
	srv.Get("/robots.txt").Do().Status(http.StatusOK)
	srv.Get("/posts/p1.html").Do().Status(http.StatusOK)
	srv.Get("/posts/not-exists.html").Do().Status(http.StatusNotFound)
	srv.Get("/themes/default/style.css").Do().Status(http.StatusOK)

	// index.html
	srv.Get("/").Do().Status(http.StatusOK)
	srv.Get("/index" + vars.Ext).Do().Status(http.StatusOK)
}
