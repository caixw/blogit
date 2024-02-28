// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serve

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v4"
	"github.com/issue9/term/v3/colors"
	"golang.org/x/text/language"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/vars"
)

func (o *options) close() error {
	return o.srv.Close()
}

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t, false)

	p, err := console.NewPrinter(language.SimplifiedChinese)
	a.NotError(err).NotNil(p)

	// 都采用默认值
	o := &options{p: p}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.source, "./").
		Equal(o.dest, "").
		Equal(o.addr, ":80")

	o = &options{
		p:    p,
		cert: "./cert",
		key:  "./key",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.addr, ":443")
}

func TestOptions_serve(t *testing.T) {
	a := assert.New(t, false)

	p, err := console.NewPrinter(language.SimplifiedChinese)
	a.NotError(err).NotNil(p)

	o := &options{
		p:      p,
		source: "../../testdata",
		addr:   ":8081",
	}

	succ := &console.Logger{Colorize: colors.New(os.Stdout)}
	info := &console.Logger{Colorize: colors.New(os.Stdout)}
	erro := &console.Logger{Colorize: colors.New(os.Stderr)}

	exit := make(chan struct{}, 1)
	go func() {
		a.Equal(o.serve(succ, info, erro), http.ErrServerClosed)
		exit <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // 等待启动完成

	// /index.html
	resp, err := http.Get("http://localhost:8081/index" + vars.Ext)
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusOK)

	// /
	resp, err = http.Get("http://localhost:8081/")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusOK)

	// /themes/default/
	resp, err = http.Get("http://localhost:8081/themes/default/")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusNotFound)

	// not-exists.html
	resp, err = http.Get("http://localhost:8081/not-exists.html")
	a.NotError(err).NotNil(resp).Equal(resp.StatusCode, http.StatusNotFound)

	a.NotError(o.close())
	<-exit
}
