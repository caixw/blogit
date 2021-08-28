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
	"sync"

	"github.com/issue9/errwrap"
	"github.com/issue9/sliceutil"

	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/loader"
	"github.com/caixw/blogit/v2/internal/vars"
)

// ErrBuilding 另一个 Rebuild 正在执行
//
// 当多次快速调用 Builder.Rebuild 时，可能返回此值，
// 表示另一个调用还未返回，新的调用又开始。
var ErrBuilding = errors.New("正在编译中")

// Builder 提供了一个可重复生成 HTML 内容的对象
type Builder struct {
	src        fs.FS
	dest       WritableFS
	rebuildMux *sync.Mutex // 防止多次调用 Rebuild

	// 以下内容在 Rebuild 之后会重新生成

	site     *site
	tpl      *template.Template
	building bool
	info     *log.Logger
}

// New 声明 Builder 实例
//
// src 需要编译的源码目录；
// dest 表示用于保存编译后的 HTML 文件的系统。可以是内存或是文件系统，
// 以及任何实现了 WritableFS 接口都可以；
func New(src fs.FS, dest WritableFS) *Builder {
	return &Builder{
		src:        src,
		dest:       dest,
		rebuildMux: &sync.Mutex{},
	}
}

// Rebuild 重新生成数据
//
// info 在运行过程中的一些提示信息通过此输出，如果为空，则会将内容写入到 io.Discard；
// base 如果不为空，则会替换 conf.yaml 配置项中的 url 字段，在预览模式下，该配置项经常会被需要修改；
func (b *Builder) Rebuild(info *log.Logger, base string) error {
	b.rebuildMux.Lock()
	defer b.rebuildMux.Unlock()

	if b.building {
		return ErrBuilding
	}

	defer func() { b.building = false }()
	b.building = true

	if info == nil {
		info = log.New(io.Discard, "", 0)
	}
	b.info = info

	if err := b.dest.Reset(); err != nil {
		return err
	}

	paths := make([]string, 0, 100)
	err := fs.WalkDir(b.src, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && !isIgnore(path) {
			paths = append(paths, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	for _, p := range paths {
		bs, err := fs.ReadFile(b.src, p)
		if err != nil {
			return err
		}
		if err = b.appendFile(loader.Slug(p), bs); err != nil {
			return err
		}
	}

	return b.buildData(base)
}

func (b *Builder) buildData(base string) (err error) {
	d, err := data.Load(b.src, base)
	if err != nil {
		return err
	}

	b.tpl, err = newTemplate(d, b.src)
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

	ignoreExts = []string{
		vars.MarkdownExt,
		".yaml", ".yml",
		".gitignore", ".git",
	}
)

func isIgnore(src string) bool {
	if ok, _ := path.Match(layoutPattern, src); ok {
		return true
	}

	if ok, _ := path.Match(themePattern, src); ok {
		return false
	}

	ext := strings.ToLower(path.Ext(src))
	return sliceutil.Index(ignoreExts, func(i int) bool { return ignoreExts[i] == ext }) >= 0
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
	b.info.Println(" >>", p)
	return b.dest.WriteFile(p, data, fs.ModePerm)
}

// Handler 将当前对象转换成 http.Handler 接口对象
//
// erro 在出错时日志的输出通道，可以为空，表示输出到 log.Default()；
func (b *Builder) Handler(erro *log.Logger) http.Handler {
	if erro == nil {
		erro = log.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: 为了自定义 index 的功能，没有采用 http.ServeFile 方法

		p := r.URL.Path
		if p != "" && p[0] == '/' {
			p = p[1:]
		}
		if p == "" || p[len(p)-1] == '/' {
			p += vars.IndexFilename
		}

		f, err := b.dest.Open(p)
		if errors.Is(err, fs.ErrNotExist) {
			http.NotFound(w, r)
			return
		} else if errors.Is(err, fs.ErrPermission) {
			errStatus(w, http.StatusForbidden)
			return
		} else if err != nil {
			erro.Println(err)
			errStatus(w, http.StatusInternalServerError)
			return
		}

		stat, err := f.Stat()
		if err != nil {
			erro.Println(err)
			errStatus(w, http.StatusInternalServerError)
			return
		}

		bs, err := io.ReadAll(f)
		if err != nil {
			erro.Println(err)
			errStatus(w, http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, p, stat.ModTime(), bytes.NewReader(bs))
	})
}

func errStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
