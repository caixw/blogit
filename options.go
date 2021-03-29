// SPDX-License-Identifier: MIT

package blogit

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/builder"
)

// Options 启动服务的参数选项
type Options struct {
	// 项目的源码目录
	//
	// 如果为空，则会采用 ./ 作为默认值。
	Src string

	// 项目编译后的输出地址
	//
	// 如果项目不大，可以采用内存文件系统，即 filesystem.Memory。
	// 你也可能根据自己的需求实现其它的，比如将内容保存到 redis 等，
	// 只要实现 filesystem.WritableFS 接口就行。
	//
	// 如果为空，则会要用 filesystem.Memory() 作为默认值。
	Dest filesystem.WritableFS

	// 如果指定了此值，那么表示要替换 conf.yaml 中的 URL
	BaseURL string

	// 服务要监听的地址
	//
	// 如果未指定此值，那么会从 BaseURL 获取相应的值。
	Addr string

	// 服务的访问根路径
	//
	// 如果此值为空，且 BaseURL 不为空，那么会从 BaseURL 获取其 Path 部分。
	Path string

	// HTTPS 模式下的证书
	//
	// 如果指定了这两个参数，将以 HTTPS 模式启动服务，否则就是普通的 HTTP 模式。
	Cert string
	Key  string

	// 不同类型日志的输出通道
	//
	// 如果为空，表示所有信息都不会输出。
	Info *log.Logger
	Erro *log.Logger
	Succ *log.Logger

	inited bool
	b      *builder.Builder
	srv    *http.Server
}

func (o *Options) serve() error {
	if o.Cert != "" && o.Key != "" {
		return o.srv.ListenAndServeTLS(o.Cert, o.Key)
	}
	return o.srv.ListenAndServe()
}

func (o *Options) init() error {
	if o.inited {
		return nil
	}

	if o.Dest == nil {
		o.Dest = filesystem.Memory()
	}

	if o.Src == "" {
		o.Src = "./"
	}

	u, err := url.Parse(o.BaseURL)
	if err != nil {
		return err
	}

	if o.Addr == "" && o.BaseURL != "" {
		o.Addr = u.Port()
		if o.Addr == "" {
			switch scheme := strings.ToLower(u.Scheme); scheme {
			case "https":
				o.Addr = ":443"
			case "http", "":
				o.Addr = ":http"
			default:
				return fmt.Errorf("不支持协议：%s", scheme)
			}
		} else {
			o.Addr = ":" + o.Addr
		}
	} else if o.Addr == "" && (o.Cert != "" && o.Key != "") {
		o.Addr = ":443"
	} else if o.Addr == "" {
		o.Addr = ":http"
	}

	if o.Path == "" {
		o.Path = u.Path
		if o.Path == "" || o.Path[0] != '/' {
			o.Path = "/" + o.Path
		}
	}

	o.b = builder.New(o.Dest, o.Erro)
	o.srv = &http.Server{Addr: o.Addr, Handler: o.initServer()}
	o.inited = true
	return nil
}

func (o *Options) initServer() http.Handler {
	o.info("启动服务：", o.Addr)

	var h http.Handler = o.b
	if o.Info != nil {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			o.Info.Printf("访问 %s\n", r.URL.String())
			o.b.ServeHTTP(w, r)
		})
	}
	return http.StripPrefix(o.Path, h)
}

func (o *Options) build() error {
	return o.b.Build(o.Src, o.BaseURL)
}

func (o *Options) info(v ...interface{}) {
	if o.Info != nil {
		o.Info.Println(v...)
	}
}

func (o *Options) erro(v ...interface{}) {
	if o.Erro != nil {
		o.Erro.Println(v...)
	}
}

func (o *Options) succ(v ...interface{}) {
	if o.Succ != nil {
		o.Succ.Println(v...)
	}
}
