// SPDX-License-Identifier: MIT

package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/blogit/internal/vars"
)

func (o *options) close() error {
	return o.srv.Close()
}

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t)

	// 都采用默认值
	o := &options{}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.source, "./").
		Equal(o.dest, "").
		Equal(o.addr, ":80")

	o = &options{
		cert: "./cert",
		key:  "./key",
	}
	a.NotError(o.sanitize()).
		Equal(o.path, "/").
		Equal(o.addr, ":443")
}

func TestOptions_serve(t *testing.T) {
	a := assert.New(t)

	o := &options{
		source: "../../../testdata/src",
		addr:   ":8080",
	}

	exit := make(chan struct{}, 1)
	go func() {
		a.Equal(o.serve(), http.ErrServerClosed)
		exit <- struct{}{}
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
