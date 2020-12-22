// SPDX-License-Identifier: MIT

package xmlwriter

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

func TestWriter_WritePI(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := New()

		w.WritePI(name, kv)
		bs := w.Bytes()
		a.Equal(string(bs), want)
	}

	test("xml-stylesheet", nil, xml.Header+"<?xml-stylesheet?>\n")
	test("xml-stylesheet", map[string]string{"type": "text/xsl"}, xml.Header+`<?xml-stylesheet type="text/xsl"?>`+"\n")
}

func TestWriter_WriteElement(t *testing.T) {
	a := assert.New(t)
	test := func(name, val string, kv map[string]string, want string) {
		w := New()

		w.WriteElement(name, val, kv)
		bs := w.Bytes()
		a.Equal(string(bs), want)
	}

	test("xml", "text", nil, xml.Header+`<xml>text</xml>`+"\n")
	test("xml", "", nil, xml.Header+`<xml></xml>`+"\n")
	test("xml", "text", map[string]string{"type": "text/xsl", "rel": "xsl"}, xml.Header+`<xml rel="xsl" type="text/xsl">text</xml>`+"\n")
}

func TestWriter_WriteCloseElement(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := New()

		w.WriteCloseElement(name, kv)
		bs := w.Bytes()
		a.Equal(string(bs), want)
	}

	test("xml", nil, xml.Header+`<xml />`+"\n")
	test("xml", nil, xml.Header+`<xml />`+"\n")
	test("xml", map[string]string{"type": "text/xsl"}, xml.Header+`<xml type="text/xsl" />`+"\n")
}
