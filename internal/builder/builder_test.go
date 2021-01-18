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

	a.NotError(os.RemoveAll("../testdata/index.xml"))

	err := Build("../testdata", "../testdata", "")
	a.NotError(err)
	a.FileExists("../testdata/index.xml")
}
