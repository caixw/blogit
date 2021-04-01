// SPDX-License-Identifier: MIT

package preview

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/vars"
)

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t)

	// 都采用默认值
	o := &options{}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.source, "./").NotNil(o.srcFS).
		Equal(o.dest, "").NotNil(o.destFS).
		Equal(o.addr, ":80")

	o = &options{
		url: "https://localhost:8080/path/",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":8080")

	o = &options{
		url: "https://localhost/path/",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":443")

	// 有证书，也有 url
	o = &options{
		url:  "http://localhost/path/",
		cert: "./cert",
		key:  "./key",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":80")

	o = &options{
		url: "ftp://localhost/path/",
	}
	a.ErrorString(o.sanitize(), "不支持")

	// url 格式错误
	o = &options{
		url: "http://localh%2%ost/path/",
	}
	a.Error(o.sanitize())
}

func TestOptions_watch(t *testing.T) {
	a := assert.New(t)

	succ := &console.Logger{Out: os.Stdout}
	info := &console.Logger{Out: os.Stdout}
	erro := &console.Logger{Out: os.Stderr}

	o := &options{
		source: "../../testdata",
		url:    "http://localhost:8080",
	}

	exit := make(chan bool, 1)
	go func() {
		a.Equal(o.watch(succ, info, erro), http.ErrServerClosed)
		exit <- true
	}()
	time.Sleep(500 * time.Millisecond) // 等待启动完成

	// /index.html
	rest.NewRequest(a, nil, http.MethodGet, "http://localhost:8080/index"+vars.Ext).
		Do().
		Status(http.StatusOK)

	// /
	rest.NewRequest(a, nil, http.MethodGet, "http://localhost:8080/").
		Do().
		Status(http.StatusOK)

	// /themes/default/
	rest.NewRequest(a, nil, http.MethodGet, "http://localhost:8080/themes/default/").
		Do().
		Status(http.StatusNotFound)

	// not-exists.html
	rest.NewRequest(a, nil, http.MethodGet, "http://localhost:8080/not-exists.html").
		Do().
		Status(http.StatusNotFound)

	a.NotError(o.close())
	<-exit
}
