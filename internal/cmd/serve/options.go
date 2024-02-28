// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package serve

import (
	"net/http"
	"os"

	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/cmd/console"
)

// 启动服务的参数选项
type options struct {
	p *message.Printer

	// 项目的源码目录
	// 如果为空，采用 ./ 作为默认值。
	source string

	// 项目编译后的输出地址
	// 如果为空，则会要用 filesystem.Memory() 作为默认值。
	dest string

	// 服务要监听的地址
	addr string

	// 服务的访问根路径
	path string

	// HTTPS 模式下的证书
	cert string
	key  string

	hookMethod string
	hookURL    string
	hookAuth   string

	b   *blogit.Builder
	srv *http.Server
}

func (o *options) serve(succ, info, erro *console.Logger) error {
	if err := o.sanitize(); err != nil {
		return err
	}

	var dest blogit.WritableFS
	if o.dest == "" {
		dest = blogit.MemoryFS()
	} else {
		dest = blogit.DirFS(o.dest)
	}
	src := os.DirFS(o.source)

	o.b = &blogit.Builder{
		Src:  src,
		Dest: dest,
		Info: info.AsLogger(),
	}
	if err := o.b.Rebuild(); err != nil {
		return err
	}

	mux := http.NewServeMux()
	h := console.Visiter(o.b.Handler(erro.AsLogger()), o.p, succ, erro)
	mux.Handle(o.path, http.StripPrefix(o.path, h))
	if o.hookURL != "" {
		mux.Handle(o.hookURL, http.HandlerFunc(o.webhook))
	}
	o.srv = &http.Server{Addr: o.addr, Handler: mux}

	info.Println(localeutil.Phrase("start server %s", o.addr).LocaleString(o.p))
	if o.cert != "" && o.key != "" {
		return o.srv.ListenAndServeTLS(o.cert, o.key)
	}
	return o.srv.ListenAndServe()
}

func (o *options) webhook(w http.ResponseWriter, r *http.Request) {
	if o.hookMethod != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if o.hookAuth != r.Header.Get("Authorization") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	o.b.Rebuild()
}

func (o *options) sanitize() error {
	if o.source == "" {
		o.source = "./"
	}

	if o.addr == "" {
		if o.cert != "" && o.key != "" {
			o.addr = ":443"
		} else {
			o.addr = ":80"
		}
	}

	if o.path == "" {
		o.path = "/"
	}

	if o.hookURL != "" {
		if o.hookAuth == "" {
			return localeutil.Error("serve hook url can not be empty")
		}

		if o.hookMethod == "" {
			return localeutil.Error("serve hook method can not be empty")
		}
	}

	return nil
}
