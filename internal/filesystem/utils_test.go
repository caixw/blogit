// SPDX-License-Identifier: MIT

package filesystem

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestExists(t *testing.T) {
	a := assert.New(t)
	fs := os.DirFS("./")

	a.True(Exists(fs, "utils.go"))
	a.True(Exists(fs, "."))
	a.False(Exists(fs, "..")) // 不允许的值
	a.False(Exists(fs, "./not-exists"))
}
