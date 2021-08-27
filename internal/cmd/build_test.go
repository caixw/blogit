// SPDX-License-Identifier: MIT

package cmd

import (
	"os"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/v2/internal/filesystem"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestCmd_Build(t *testing.T) {
	a := assert.New(t)
	dest, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)

	a.NotError(Exec([]string{"build", "-src", "../testdata", "-dest", dest}))

	fs := os.DirFS(dest)
	a.True(filesystem.Exists(fs, "index"+vars.Ext))
}
