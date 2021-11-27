// SPDX-License-Identifier: MIT

package builder

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestStripTags(t *testing.T) {
	a := assert.New(t, false)

	tests := map[string]string{
		"<div>str</div>":        "str",
		"str<br />":             "str",
		"<div><p>str</p></div>": "str",
	}

	for expr, val := range tests {
		a.Equal(stripTags(expr), val, "测试[%v]时出错", expr)
	}
}
