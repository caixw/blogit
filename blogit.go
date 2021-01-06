// SPDX-License-Identifier: MIT

// Package blogit 依赖于 git 的博客系统
package blogit

import (
	"net/http"

	"github.com/caixw/blogit/internal/builder"
	"github.com/caixw/blogit/internal/data"
)

// Version 版本号
const Version = "0.1.0"

// Build 编译并输出内容
func Build(dir string) error {
	b, err := newBuilder(dir)
	if err != nil {
		return err
	}
	return b.Dump(dir)
}

// Serve 运行服务
func Serve(src, addr, path string) error {
	b, err := Handler(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
func ServeTLS(src, addr, path, cert, key string) error {
	b, err := Handler(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServeTLS(addr, cert, key, nil)
}

// Handler 将编译后的内容作为 http.Handler 接口返回
func Handler(dir string) (http.Handler, error) {
	return newBuilder(dir)
}

func newBuilder(dir string) (*builder.Builder, error) {
	d, err := data.Load(dir)
	if err != nil {
		return nil, err
	}

	b := &builder.Builder{}
	if err := b.Load(d); err != nil {
		return nil, err
	}
	return b, nil
}
