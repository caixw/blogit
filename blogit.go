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
func Build(src, dest string) error {
	return builder.Build(src, dest)
}

// Serve 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func Serve(src, addr, path string, l *log.Logger) error {
	return ServeTLS(src, addr, path, "", "", l)
}

// ServeTLS 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func ServeTLS(src, addr, path, cert, key string, l *log.Logger) error {
	b := &builder.Builder{}
	if err := b.Build(src, ""); err != nil {
		return err
	}
	return serve(b, src, addr, path, cert, key, l)
}

func serve(b *builder.Builder, src, addr, path, cert, key string, l *log.Logger) error {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}

	var h http.Handler = b

	if l != nil {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Printf("访问 %s\n", r.URL.String())
			b.ServeHTTP(w, r)
		})
	}

	http.Handle(path, http.StripPrefix(path, h))

	if cert != "" {
		return http.ListenAndServeTLS(addr, cert, key, nil)
	}
	return http.ListenAndServe(addr, nil)
}
