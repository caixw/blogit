// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLicense_sanitize(t *testing.T) {
	a := assert.New(t)

	l := &License{}
	a.Error(l.sanitize())

	l = &License{Text: "t"}
	a.Error(l.sanitize())

	l = &License{Text: "t", URL: "/favicon.ico"}
	a.NotError(l.sanitize())
}

func TestMenu_sanitize(t *testing.T) {
	a := assert.New(t)

	m := &Menu{}
	a.Error(m.sanitize())

	m = &Menu{Text: "t"}
	a.Error(m.sanitize())

	m = &Menu{Text: "t", URL: "/favicon.ico"}
	a.NotError(m.sanitize())
}

func TestAuthor_sanitize(t *testing.T) {
	a := assert.New(t)

	author := &Author{}
	a.Error(author.sanitize())

	author.Name = ""
	a.Error(author.sanitize())

	author.Name = "caixw"
	a.NotError(author.sanitize())

	author.Email = "invalid-email"
	a.Error(author.sanitize())
}

func TestIcon_sanitize(t *testing.T) {
	a := assert.New(t)

	icon := &Icon{}
	a.Error(icon.sanitize())

	// type 会根据 URL 的扩展名，可能会自动计算获得。
	// png 有内置，肯定可以成功检测到
	icon = &Icon{URL: "http://example.com/1.png"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "image/png")

	// 不存在的扩展名，则不会计算工其 type
	icon = &Icon{URL: "http://example.com/1.not-exists"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "")
}
