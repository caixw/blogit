// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
package blogit

import (
	"log"
	"net/http"

	"github.com/caixw/blogit/internal/builder"
)

// Version 版本号
const Version = "0.1.0"

// Build 编译并输出内容
func Build(dir string) error {
	return builder.Build(dir)
}

// Serve 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func Serve(src, addr, path string, l *log.Logger) error {
	http.Handle(path, newHandler(path, src, l))
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
//
// 如果 l 不为 nil，则会在此通道上输出访问记录；
func ServeTLS(src, addr, path, cert, key string, l *log.Logger) error {
	http.Handle(path, newHandler(path, src, l))
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
