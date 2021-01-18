// SPDX-License-Identifier: MIT

package builder

import (
	"os"
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestFT(t *testing.T) {
	a := assert.New(t)

	a.Empty(ft(time.Time{}))
	a.NotEmpty(ft(time.Now()))
}

func TestNewHTML(t *testing.T) {
	a := assert.New(t)

	a.Nil(newHTML(""))
	a.NotNil(newHTML(" "))
}

func TestBuild(t *testing.T) {
	a := assert.New(t)

	a.NotError(os.RemoveAll("../../testdata/dest/index.xml"))

	err := Build("../../testdata/src", "../../testdata/dest", "")
	a.NotError(err)
	a.FileExists("../../testdata/dest/index.xml")
}
