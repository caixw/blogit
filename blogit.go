// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
package blogit

import (
	"log"
	"net/http"

	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/vars"
)

// Version 版本号
func Version() string {
	return vars.Version()
}

// Build 编译并输出内容
//
// dir 表示源码目录；
// dest 表示输出的目录；
// base 表示网站的基地址，如果此值不为空，会替代 conf.yaml 中的 url 变量，
// 在预览模式下，此参数会很有用。
func Build(src, dest, base string) error {
	return builder.Build(src, dest, base)
}

// Serve 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func Serve(dest, addr, path string, l *log.Logger) error {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	http.Handle(path, newHandler(path, dest, l))
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func ServeTLS(dest, addr, path, cert, key string, l *log.Logger) error {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	http.Handle(path, newHandler(path, dest, l))
	return http.ListenAndServeTLS(addr, cert, key, nil)
}

func newHandler(prefix, dir string, l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l != nil {
			l.Printf("访问 %s\n", r.URL.String())
		}

		http.StripPrefix(prefix, http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	})
}
