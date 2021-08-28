// SPDX-License-Identifier: MIT

package serve

import (
	"net/http"
	"os"

	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/builder"
	"github.com/caixw/blogit/v2/internal/cmd/console"
)

// options 启动服务的参数选项
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

	srv *http.Server
}

func (o *options) serve(info, erro *console.Logger) error {
	if err := o.sanitize(); err != nil {
		return err
	}

	var dest builder.WritableFS
	if o.dest == "" {
		dest = builder.MemoryFS()
	} else {
		dest = builder.DirFS(o.dest)
	}
	src := os.DirFS(o.source)

	b := blogit.NewBuilder(src, dest)
	if err := b.Rebuild(info.AsLogger(), ""); err != nil {
		return err
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info.Println(o.p.Sprintf("visit url", r.URL.String()))
		b.Handler(erro.AsLogger()).ServeHTTP(w, r)
	})
	o.srv = &http.Server{Addr: o.addr, Handler: http.StripPrefix(o.path, h)}

	info.Println(o.p.Sprintf("start server", o.addr))
	if o.cert != "" && o.key != "" {
		return o.srv.ListenAndServeTLS(o.cert, o.key)
	}
	return o.srv.ListenAndServe()
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

	return nil
}
