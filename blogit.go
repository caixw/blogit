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
func Build(src, target string) error {
	b, err := newBuilder(src)
	if err != nil {
		return err
	}
	return b.Dump(target)
}

// Serve 运行服务
func Serve(src, addr, path string) error {
	b, err := newBuilder(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServe(addr, nil)
}

// ServeTLS 运行服务
func ServeTLS(src, addr, path, cert, key string) error {
	b, err := newBuilder(src)
	if err != nil {
		return err
	}

	http.Handle(path, b)
	return http.ListenAndServeTLS(addr, cert, key, nil)
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
