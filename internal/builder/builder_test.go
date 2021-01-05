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

func newBuilder(a *assert.Assertion, dir string) *Builder {
	d, err := data.Load(dir)
	a.NotError(err).NotNil(d)

	b, err := Build(d)
	a.NotError(err).NotNil(b)

	return b
}

func TestBuild(t *testing.T) {
	a := assert.New(t)

	b := newBuilder(a, "../testdata")
	a.Equal(b.Builded.Year(), time.Now().Year())
}
