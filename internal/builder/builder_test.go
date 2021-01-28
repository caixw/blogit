// SPDX-License-Identifier: MIT

package builder

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

func TestBuild(t *testing.T) {
	a := assert.New(t)

	a.NotError(os.RemoveAll("../../testdata/dest/index.xml"))

	err := Build("../../testdata/src", "../../testdata/dest", "")
	a.NotError(err)
	a.FileExists("../../testdata/dest/index.html")
}
