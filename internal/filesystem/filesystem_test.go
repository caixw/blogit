// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package filesystem

import (
	"os"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestExists(t *testing.T) {
	a := assert.New(t, false)
	fs := os.DirFS("./")

	a.True(Exists(fs, "filesystem.go"))
	a.True(Exists(fs, "."))
	a.False(Exists(fs, "..")) // 不允许的值
	a.False(Exists(fs, "./not-exists"))
}
