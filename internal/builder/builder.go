// SPDX-License-Identifier: MIT

package builder

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/issue9/errwrap"
	"github.com/otiai10/copy"

	"github.com/caixw/blogit/internal/data"
	"github.com/caixw/blogit/internal/vars"
)

var copyOptions = copy.Options{
	Skip: func(src string) (bool, error) {
		// 忽略 themes/xxx/layout
		if filepath.Base(filepath.Dir(src)) == vars.LayoutDir || filepath.Base(src) == vars.LayoutDir {
			return true, nil
		}

		ext := strings.ToLower(filepath.Ext(src))
		return ext == vars.MarkdownExt ||
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
	site  *site
	tpl   *template.Template
	files map[string][]byte
}

// Build 编译内容
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
	d, err := data.Load(dir, base)
	if err != nil {
		return nil, err
	}

	tpl, err := newTemplate(d, dir)
	if err != nil {
		return nil, err
	}
	b := &builder{
		site:  newSite(d),
		tpl:   tpl,
		files: make(map[string][]byte, 20),
	}

	if err := b.buildTags(d); err != nil {
		return nil, err
	}

	if err := b.buildPosts(d); err != nil {
		return nil, err
	}

	if err := b.buildSitemap(d); err != nil {
		return nil, err
	}

	if err := b.buildArchive(d); err != nil {
		return nil, err
	}

	if err := b.buildAtom(d); err != nil {
		return nil, err
	}

	if err := b.buildRSS(d); err != nil {
		return nil, err
	}

	if err := b.buildRobots(d); err != nil {
		return nil, err
	}

	if err := b.buildProfile(d); err != nil {
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
func (b *builder) appendTemplateFile(path string, p *page) error {
	buf := &bytes.Buffer{}

	if err := b.tpl.ExecuteTemplate(buf, p.Type, p); err != nil {
		return err
	}

	b.files[path] = buf.Bytes()
	return nil
}

// path 表示输出的文件路径，相对于源目录；
// xsl 表示关联的 xsl，如果不需要则可能为空；
func (b *builder) appendXMLFile(d *data.Data, path, xsl string, v interface{}) error {
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

	b.files[path] = buf.Bytes()
	return nil
}
