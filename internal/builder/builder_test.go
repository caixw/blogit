// SPDX-License-Identifier: MIT

package builder

import (
	"net/http"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/internal/data"
)

var _ http.Handler = &Builder{}

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

func newBuilder(a *assert.Assertion, dir string) *Builder {
	d, err := data.Load(dir)
	a.NotError(err).NotNil(d)

	b := &Builder{}
	err = b.Load(d)
	a.NotError(err).NotNil(b)

	return b
}

func TestBuild(t *testing.T) {
	a := assert.New(t)

	b := newBuilder(a, "../testdata")
	a.Equal(b.Builded.Year(), time.Now().Year())
}
