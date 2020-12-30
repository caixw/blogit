// SPDX-License-Identifier: MIT

package builder

import (
	"testing"

	"github.com/issue9/assert"
)

func TestData_BuildURL(t *testing.T) {
	a := assert.New(t)

	data := &Data{URL: "https://example.com/"} // 传入的 URL 必定是以 / 结尾的
	a.Equal(data.BuildURL("/p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(data.BuildURL("p1/p2.md"), "https://example.com/p1/p2.md")
	a.Equal(data.BuildURL(""), "https://example.com/")
	a.Equal(data.BuildURL("/"), "https://example.com/")
}
