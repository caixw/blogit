// SPDX-License-Identifier: MIT

package preview

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
	"github.com/issue9/term/v3/colors"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t, false)

	p, err := console.NewPrinter()
	a.NotError(err).NotNil(p)

	// 都采用默认值
	o := &options{p: p}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.source, "./").NotNil(o.srcFS).
		Equal(o.dest, "").NotNil(o.destFS).
		Equal(o.addr, ":80")

	o = &options{
		p:   p,
		url: "https://localhost:8080/path/",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":8080")

	o = &options{
		p:   p,
		url: "https://localhost/path/",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":443")

	// 有证书，也有 url
	o = &options{
		p:    p,
		url:  "http://localhost/path/",
		cert: "./cert",
		key:  "./key",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/path/").
		Equal(o.addr, ":80")

	o = &options{
		p:   p,
		url: "ftp://localhost/path/",
	}
	a.ErrorString(o.sanitize(), "未支持的协议")

	// url 格式错误
	o = &options{
		p:   p,
		url: "http://localh%2%ost/path/",
	}
	a.Error(o.sanitize())
}

func TestOptions_watch(t *testing.T) {
	a := assert.New(t, false)

	p, err := console.NewPrinter()
	a.NotError(err).NotNil(p)

	succ := &console.Logger{Colorize: colors.New(os.Stdout)}
	info := &console.Logger{Colorize: colors.New(os.Stdout)}
	erro := &console.Logger{Colorize: colors.New(os.Stderr)}

	o := &options{
		p:      p,
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
	resp, err := http.Get("http://localhost:8080/index" + vars.Ext)
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusOK)

	// /
	resp, err = http.Get("http://localhost:8080/")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusOK)

	// /themes/default/
	resp, err = http.Get("http://localhost:8080/themes/default/")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusNotFound)

	// not-exists.html
	resp, err = http.Get("http://localhost:8080/not-exists.html")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusNotFound)

	a.NotError(o.close())
	<-exit
}
