// SPDX-License-Identifier: MIT

// Package builder 提供编译成 HTML 的相关功能
package builder

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/utils"
	"github.com/caixw/blogit/internal/vars"
)

type file struct {
	data []byte
	ct   string
	path string
}

// Builder 提供了一个可重复生成 HTML 内容的对象
type Builder struct {
	site  *site
	tpl   *template.Template
	files []*file
}

// Build 重新生成数据
func (b *Builder) Build(src, base string) error {
	paths := make([]string, 0, 100)

	err := filepath.Walk(src, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && !isIgnore(path) {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if b.files == nil {
		b.files = make([]*file, 0, len(paths))
	} else {
		b.files = b.files[:0]
	}

	dir := filepath.ToSlash(filepath.Clean(src))

	for _, p := range paths {
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		p = filepath.ToSlash(filepath.Clean(p))
		p = strings.TrimPrefix(p, dir)

		b.appendFile(p, "", data)
	}

	return b.buildData(src, base)
}

func (b *Builder) buildData(src, base string) (err error) {
	d, err := data.Load(src, base)
	if err != nil {
		return err
	}

	b.tpl, err = newTemplate(d, src)
	if err != nil {
		return err
	}

	b.site = newSite(d)

	call := func(f func(*data.Data) error) {
		if err == nil {
			err = f(d)
		}
	}

	call(b.buildTags)
	call(b.buildPosts)
	call(b.buildSitemap)
	call(b.buildArchive)
	call(b.buildAtom)
	call(b.buildRSS)
	call(b.buildRobots)
	call(b.buildProfile)

	return
}

func isIgnore(src string) bool {
	// themes/**/layout/file 这种格式将忽略
	layout := path.Dir(src)
	if path.Base(layout) == vars.LayoutDir &&
		path.Base(path.Dir(path.Dir(layout))) == vars.ThemesDir {
		return true
	}

	ext := strings.ToLower(path.Ext(src))
	return ext == vars.MarkdownExt ||
		ext == ".yaml" ||
		ext == ".yml" ||
		ext == ".gitignore" ||
		ext == ".git"
}

// Build 编译内容
func Build(src, dest string) error {
	b := &Builder{}
	if err := b.Build(src, ""); err != nil {
		return err
	}
	return b.Dump(dest)
}

// Dump 将内容输出到 dir 目录
func (b *Builder) Dump(dir string) error {
	if err := os.MkdirAll(filepath.Join(dir, vars.TagsDir), os.ModePerm); err != nil {
		return err
	}

	for _, f := range b.files {
		path := filepath.Join(dir, f.path)

		base := filepath.Dir(path)
		if !utils.FileExists(base) {
			if err := os.MkdirAll(base, os.ModePerm); err != nil {
				return err
			}
		}

		err := ioutil.WriteFile(path, f.data, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// path 表示输出的文件路径，相对于源目录；
func (b *Builder) appendTemplateFile(path string, p *page) error {
	buf := &bytes.Buffer{}

	if err := b.tpl.ExecuteTemplate(buf, p.Type, p); err != nil {
		return err
	}

	b.appendFile(path, "", buf.Bytes())
	return nil
}

// path 表示输出的文件路径，相对于源目录；
// xsl 表示关联的 xsl，如果不需要则可能为空；
func (b *Builder) appendXMLFile(d *data.Data, path, xsl string, v interface{}) error {
	data, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	buf := &errwrap.Buffer{}
	buf.WString(xml.Header)
	if xsl != "" {
		buf.Printf(`<?xml-stylesheet type="text/xsl" href="%s"?>`, xsl).WByte('\n')
	}
	buf.WBytes(data)

	if buf.Err != nil {
		return buf.Err
	}

	b.appendFile(path, "application/xml", buf.Bytes())
	return nil
}

// 如果 path 以 / 开头，则会自动去除 /
func (b *Builder) appendFile(p, ct string, data []byte) {
	if p == "" {
		panic("参数 path 不能为空")
	}
	if p[0] == '/' {
		p = p[1:]
	}

	if ct == "" {
		ct = mime.TypeByExtension(path.Ext(p))
	}
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	if index := strings.IndexByte(ct, ';'); index > 0 {
		ct = ct[:index]
	}

	b.files = append(b.files, &file{
		data: data,
		path: p,
		ct:   ct,
	})
}

// ServeHTTP 作为 HTTP 服务接口使用
func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p != "" && p[0] == '/' {
		p = p[1:]
	}

	if p == "" {
		p = "index" + vars.Ext
	}

	for _, f := range b.files {
		if f.path == p {
			w.Header().Set("Content-Type", f.ct)
			w.Write(f.data)
			return
		}
	}

	http.NotFound(w, r)
}
