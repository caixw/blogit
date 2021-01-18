// SPDX-License-Identifier: MIT

package builder

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/issue9/errwrap"
	"github.com/otiai10/copy"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

const (
	// 输出的时间格式
	//
	// NOTE: 时间可能会被当作 XML 的属性值，如果格式中带引号，需要注意正确处理。
	timeFormat = time.RFC3339
)

var copyOptions = copy.Options{
	Skip: func(src string) (bool, error) {
		ext := strings.ToLower(filepath.Ext(src))
		return ext == ".md" ||
			ext == ".yaml" ||
			ext == ".yml" ||
			ext == ".gitignore" ||
			ext == ".git", nil
	},

	OnSymlink: func(src string) copy.SymlinkAction {
		return copy.Skip
	},
	AddPermission: 0200,
}

type builder struct {
	files map[string][]byte
}

type innerhtml struct {
	Content string `xml:",innerxml"`
}

func ft(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeFormat)
}

func newHTML(html string) *innerhtml {
	if html == "" {
		return nil
	}
	return &innerhtml{Content: html}
}

// Build 编译成 xml 文件
func Build(src, dest, base string) error {
	if src != dest {
		if err := copy.Copy(src, dest, copyOptions); err != nil {
			return err
		}
	}

	b, err := newBuilder(src, base)
	if err != nil {
		return err
	}

	return b.dump(dest)
}

func newBuilder(dir, base string) (*builder, error) {
	d, err := data.Load(dir)
	if err != nil {
		return nil, err
	}

	// base 被用于替换 data.URL，所以了要和其有一样规则：保证以 / 结尾。
	if base != "" && base[len(base)-1] != '/' {
		base += "/"
	}

	if base != "" {
		d.URL = base
	}

	b := &builder{
		files: make(map[string][]byte, 20),
	}

	if err := b.buildInfo(vars.InfoXML, d); err != nil {
		return nil, err
	}

	if err := b.buildTags(d); err != nil {
		return nil, err
	}

	if err := b.buildPosts(d); err != nil {
		return nil, err
	}

	if err := b.buildSitemap(vars.SitemapXML, d); err != nil {
		return nil, err
	}

	if err := b.buildArchive(vars.ArchiveXML, d); err != nil {
		return nil, err
	}

	if err := b.buildAtom(vars.AtomXML, d); err != nil {
		return nil, err
	}

	if err := b.buildRSS(vars.RssXML, d); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *builder) dump(dir string) error {
	if err := os.MkdirAll(filepath.Join(dir, vars.TagsDir), os.ModePerm); err != nil {
		return err
	}

	for path, content := range b.files {
		err := ioutil.WriteFile(filepath.Join(dir, path), content, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// path 表示输出的文件路径，相对于源目录；
// xsl 表示关联的 xsl，相对于当前主题目录的路径，如果不需要则可能为空；
// ct 表示内容的 content-type 值，为空表示采用 application/xml；
func (b *builder) appendXMLFile(d *data.Data, path, xsl string, lastmod time.Time, v interface{}) error {
	data, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	buf := &errwrap.Buffer{}
	buf.WString(xml.Header)
	if xsl != "" {
		xsl = d.BuildThemeURL(xsl)
		buf.Printf(`<?xml-stylesheet type="text/xsl" href="%s"?>`, xsl).WByte('\n')
	}
	buf.WBytes(data)

	if buf.Err != nil {
		return buf.Err
	}

	b.files[path] = buf.Bytes()
	return nil
}
