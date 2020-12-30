// SPDX-License-Identifier: MIT

package utils

import (
	"testing"

	"github.com/issue9/assert"
)

func TestFileExists(t *testing.T) {
	a := assert.New(t)

	a.True(FileExists("./utils.go"))
	a.True(FileExists("."))
	a.True(FileExists("../"))
	a.False(FileExists("../not-exists"))
}
