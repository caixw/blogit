// SPDX-License-Identifier: MIT

package loader

import (
	"bytes"
	"io/fs"

	fh "github.com/alecthomas/chroma/formatters/html"
	toc "github.com/mdigger/goldmark-toc"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/internal/vars"
)

var (
	markdown = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Strikethrough,
			extension.Footnote,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					fh.WithLineNumbers(true),
					fh.WithClasses(true),
					fh.ClassPrefix(vars.HighlightClassPrefix),
				),
			),
		),

		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),

		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
)

func convert(fsys fs.FS, path string) (*Post, error) {
	bs, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	ctx := parser.NewContext(parser.WithIDs(toc.NewIDs("")))
	buf := new(bytes.Buffer)

	doc := markdown.Parser().Parse(text.NewReader(bs), parser.WithContext(ctx))
	headers := toc.Headers(doc, bs)
	if err = markdown.Renderer().Render(buf, bs, doc); err != nil {
		return nil, err
	}

	metadata, err := yaml.Marshal(meta.Get(ctx))
	if err != nil {
		return nil, err
	}
	post := &Post{}
	if err := yaml.Unmarshal(metadata, post); err != nil {
		return nil, err
	}
	post.Content = buf.String()

	start := 6
	for _, h := range headers {
		if start > h.Level {
			start = h.Level
		}
	}

	toc := make([]Header, 0, len(headers))
	for _, h := range headers {
		toc = append(toc, Header{
			Indent: h.Level - start,
			Level:  h.Level,
			ID:     h.ID,
			Text:   h.Text,
		})
	}
	post.TOC = toc

	return post, nil
}
