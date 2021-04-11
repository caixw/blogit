// SPDX-License-Identifier: MIT

// Package builder 提供编译成 HTML 的相关功能
package builder

import (
	"bytes"
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

// Builder 提供了一个可重复生成 HTML 内容的对象
type Builder struct {
	info *log.Logger
	erro *log.Logger
	wfs  WritableFS

	// 以下内容在 ReBuild 之后会重新生成
	site *site
	tpl  *template.Template
}

// New 声明 Builder 实例
//
// wfs 表示用于保存编译后的 HTML 文件的系统。可以是内存或是文件系统，
// 以及任何实现了 WritableFS 接口都可以；
// info 在运行过程中的一些提示信息通过此输出，如果为空，则会采用 log.Default()；
// erro 表示的是在把 Builder 当作 http.Handler 处理时，在出错时的日志输出通道。
// 如果为空，则会采用 log.Default()。如果不准备其当作 http.Handler 使用，则此值是无用；
func New(wfs WritableFS, info, erro *log.Logger) *Builder {
	if erro == nil {
		erro = log.Default()
	}
	if info == nil {
		info = log.Default()
	}
	return &Builder{wfs: wfs, info: info, erro: erro}
}

// Rebuild 重新生成数据
func (b *Builder) Rebuild(src fs.FS, base string) error {
	if err := b.wfs.Reset(); err != nil {
		return err
	}

	paths := make([]string, 0, 100)
	err := fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && !isIgnore(path) {
			paths = append(paths, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	for _, p := range paths {
		bs, err := fs.ReadFile(src, p)
		if err != nil {
			return err
		}
		if err = b.appendFile(loader.Slug(p), bs); err != nil {
			return err
		}
	}

	return b.buildData(src, base)
}

func (b *Builder) buildData(src fs.FS, base string) (err error) {
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
	call(b.buildIndexes)
	call(b.buildSitemap)
	call(b.buildArchive)
	call(b.buildAtom)
	call(b.buildRSS)
	call(b.buildRobots)
	call(b.buildProfile)
	call(b.buildHighlights)

	return
}

var (
	layoutPattern = path.Join(vars.ThemesDir, "*", vars.LayoutDir, "*")
	themePattern  = path.Join(vars.ThemesDir, "*", vars.ThemeYAML)
)

func isIgnore(src string) bool {
	if ok, _ := path.Match(layoutPattern, src); ok {
		return true
	}

	if ok, _ := path.Match(themePattern, src); ok {
		return false
	}

	ext := strings.ToLower(path.Ext(src))
	return ext == vars.MarkdownExt ||
		ext == ".yaml" ||
		ext == ".yml" ||
		ext == ".gitignore" ||
		ext == ".git"
}

// path 表示输出的文件路径，相对于源目录；
// xsl 表示关联的 xsl，如果不需要则可能为空；
func (b *Builder) appendXMLFile(path, xsl string, v interface{}) error {
	bs, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	buf := &errwrap.Buffer{}
	buf.WString(xml.Header)
	if xsl != "" {
		buf.Printf(`<?xml-stylesheet type="text/xsl" href="%s"?>`, xsl).WByte('\n')
	}
	buf.WBytes(bs)

	if buf.Err != nil {
		return buf.Err
	}

	return b.appendFile(path, buf.Bytes())
}

// 如果 path 以 / 开头，则会自动去除 /
func (b *Builder) appendFile(p string, data []byte) error {
	b.info.Println("添加：", p)
	return b.wfs.WriteFile(p, data, fs.ModePerm)
}

// ServeHTTP 作为 HTTP 服务接口使用
func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: 为了自定义 index 的功能，没有采用 http.ServeFile 方法

	p := r.URL.Path
	if p != "" && p[0] == '/' {
		p = p[1:]
	}
	if p == "" || p[len(p)-1] == '/' {
		p += vars.IndexFilename
	}

	f, err := b.wfs.Open(p)
	if errors.Is(err, fs.ErrNotExist) {
		http.NotFound(w, r)
		return
	} else if errors.Is(err, fs.ErrPermission) {
		errStatus(w, http.StatusForbidden)
		return
	} else if err != nil {
		b.erro.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	stat, err := f.Stat()
	if err != nil {
		b.erro.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		b.erro.Println(err)
		errStatus(w, http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, p, stat.ModTime(), bytes.NewReader(bs))
}

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
