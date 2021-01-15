// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
package blogit

import (
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
func Serve(src, addr, path string) error {
	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(src))))
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
func ServeTLS(src, addr, path, cert, key string) error {
	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(src))))
	return http.ListenAndServeTLS(addr, cert, key, nil)
}
