// SPDX-License-Identifier: MIT

package blogit

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/blogit/internal/vars"
)

func TestServe(t *testing.T) {
	a := assert.New(t)

	o := &Options{
		Src:  os.DirFS("./testdata/src"),
		Addr: ":8080",
		Erro: log.Default(),
		Info: log.Default(),
		Succ: log.Default(),
	}
	s, err := Serve(o)
	a.NotError(err).NotNil(s)

	exit := make(chan struct{}, 1)
	go func() {
		a.Equal(s.Serve(), http.ErrServerClosed)
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

	a.NotError(s.Close())
	<-exit
}

func TestWatch(t *testing.T) {
	a := assert.New(t)

	o := &Options{
		Src:  os.DirFS("./testdata/src"),
		Addr: ":8080",
		Erro: log.Default(),
		Info: log.Default(),
		Succ: log.Default(),
	}
	s, err := Watch("./testdata/src", o)
	a.NotError(err).NotNil(s)

	exit := make(chan struct{}, 1)
	go func() {
		a.Equal(s.Serve(), http.ErrServerClosed)
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

	a.NotError(s.Close())
	<-exit
}
