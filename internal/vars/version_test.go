// SPDX-License-Identifier: MIT

package vars

import (
	"testing"

	"github.com/issue9/assert"
	v "github.com/issue9/version"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(v.SemVerValid(Version()))
}
