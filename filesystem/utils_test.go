// SPDX-License-Identifier: MIT

package filesystem

import (
	"testing"

	"github.com/issue9/assert"
)

func TestExists(t *testing.T) {
	a := assert.New(t)

	a.True(Exists("./utils.go"))
	a.True(Exists("."))
	a.True(Exists("../"))
	a.False(Exists("../not-exists"))
}
