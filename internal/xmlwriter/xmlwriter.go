// SPDX-License-Identifier: MIT

// Package xmlwriter XML 写入工具
package xmlwriter

import (
	"encoding/xml"
	"sort"
	"strings"

	"github.com/issue9/errwrap"
)

// XMLWriter XML 写操作
//
// 相对于官方的 XML 包，因为不涉及到反射操作，性能上会有所提升。
type XMLWriter struct {
	errwrap.Buffer
	indent int // 保存当前的缩进量
}

// New 声明一个新的 XMLWriter
func New() *XMLWriter {
	w := &XMLWriter{}
	w.WString(xml.Header)
	return w
}

func (w *XMLWriter) writeIndent() {
	w.WString(strings.Repeat(" ", w.indent*4))
}

// WriteStartElement 写入一个开始元素
func (w *XMLWriter) WriteStartElement(name string, attr map[string]string) {
	w.startElement(name, attr, true)
}

// newline 是否换行
func (w *XMLWriter) startElement(name string, attr map[string]string, newline bool) {
	w.writeIndent()
	w.indent++

	w.WByte('<')
	w.WString(name)
	w.writeAttr(attr)
	w.WByte('>')

	if newline {
		w.WByte('\n')
	}
}

// WriteEndElement 写入一个结束元素
func (w *XMLWriter) WriteEndElement(name string) {
	w.endElement(name, true)
}

// indent 是否需要填上缩进时的字符，如果不换行输出结束符，则不能输出缩进字符串
func (w *XMLWriter) endElement(name string, indent bool) {
	w.indent--
	if indent {
		w.writeIndent()
	}
	w.Printf("</%s>\n", name)
}

// WriteCloseElement 写入一个自闭合的元素
//
// name 元素标签名；
// attr 元素的属性。
func (w *XMLWriter) WriteCloseElement(name string, attr map[string]string) {
	w.writeIndent()

	w.WByte('<')
	w.WString(name)
	w.writeAttr(attr)
	w.WString(" />\n")
}

// WriteElement 写入一个完整的元素
//
// name 元素标签名；
// val 元素内容；
// attr 元素的属性。
func (w *XMLWriter) WriteElement(name, val string, attr map[string]string) {
	w.startElement(name, attr, false)
	w.WString(val)
	w.endElement(name, false)
}

// WritePI 写入一个 PI 指令
func (w *XMLWriter) WritePI(name string, kv map[string]string) {
	w.WString("<?")
	w.WString(name)
	w.writeAttr(kv)
	w.WString("?>\n")
}

func (w *XMLWriter) writeAttr(attr map[string]string) {
	keys := make([]string, 0, len(attr))
	for key := range attr {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		w.Printf(` %s="%s"`, key, attr[key])
	}
}

// Bytes 将内容转换成 []byte 并返回
func (w *XMLWriter) Bytes() []byte {
	if w.Err != nil {
		panic(w.Err)
	}
	return w.Buffer.Bytes()
}
