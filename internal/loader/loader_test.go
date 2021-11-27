// SPDX-License-Identifier: MIT

package loader

import (
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/localeutil"
)

var (
	_ localeutil.LocaleStringer = &FieldError{}
	_ error                     = &FieldError{}
)

func TestLicense_sanitize(t *testing.T) {
	a := assert.New(t, false)

	l := &Link{}
	a.Error(l.sanitize())

	l = &Link{Text: "t"}
	a.Error(l.sanitize())

	l = &Link{Text: "t", URL: "/favicon.ico"}
	a.NotError(l.sanitize())
}

func TestAuthor_sanitize(t *testing.T) {
	a := assert.New(t, false)

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
	a := assert.New(t, false)

	icon := &Icon{}
	a.Error(icon.sanitize())

	// type 会根据 URL 的扩展名，可能会自动计算获得。
	// png 有内置，肯定可以成功检测到
	icon = &Icon{URL: "https://example.com/1.png"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "image/png")

	// 不存在的扩展名，则不会计算工其 type
	icon = &Icon{URL: "https://example.com/1.not-exists"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "")
}
