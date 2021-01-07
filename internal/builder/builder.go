// SPDX-License-Identifier: MIT

package builder

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/issue9/errwrap"

	"github.com/caixw/blogit/internal/data"
)

const xmlContentType = "application/xml"

// Builder 保存构建好的数据
type Builder struct {
	files   []*file
	Builded time.Time
}

type file struct {
	path    string
	lastmod time.Time
	content []byte
	ct      string
}

type datetime struct {
	Long  string `xml:"long,attr"`
	Short string `xml:"short,attr"`
}

type innerhtml struct {
	Content string `xml:",innerxml"`
}

func toDatetime(t time.Time, d *data.Data) datetime {
	return datetime{
		Long:  t.Format(d.LongDateFormat),
		Short: t.Format(d.ShortDateFormat),
	}
}

// Load 加载数据到当前实例
func (b *Builder) Load(d *data.Data) error {
	b.files = make([]*file, 0, 20)
	b.Builded = d.Builded

	if err := b.buildInfo("info.xml", d); err != nil {
		return err
	}

	if err := b.buildTags(d); err != nil {
		return err
	}

	if err := b.buildPosts(d); err != nil {
		return err
	}

	if err := b.buildSitemap("sitemap.xml", d); err != nil {
		return err
	}

	if err := b.buildArchives("archives.xml", d); err != nil {
		return err
	}

	if err := b.buildAtom("atom.xml", d); err != nil {
		return err
	}

	if err := b.buildRSS("rss.xml", d); err != nil {
		return err
	}

	return nil
}

func (f *file) dump(dir string) error {
	return ioutil.WriteFile(filepath.Join(dir, f.path), f.content, os.ModePerm)
}

// Dump 输出内容
func (b *Builder) Dump(dir string) error {
	for _, f := range b.files {
		if err := f.dump(dir); err != nil {
			return err
		}
	}
	return nil
}

// ServeHTTP 以内容进行 HTTP 服务
func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, f := range b.files {
		if f.path == path {
			w.Header().Set("Content-Type", f.ct)
			http.ServeContent(w, r, f.path, f.lastmod, bytes.NewReader(f.content))
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// xsl 表示关联的 xsl，如果不需要则可能为空；
// ct 表示内容的 content-type 值，为空表示采用 application/xml；
func (b *Builder) appendXMLFile(path, xsl string, lastmod time.Time, v interface{}) error {
	data, err := xml.Marshal(v)
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

	b.files = append(b.files, &file{
		path:    path,
		lastmod: lastmod,
		content: buf.Bytes(),
		ct:      xmlContentType,
	})
	return nil
}
