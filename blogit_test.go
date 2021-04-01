// SPDX-License-Identifier: MIT

package blogit

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/blogit/builder"
	"github.com/caixw/blogit/internal/filesystem"
	"github.com/caixw/blogit/internal/testdata"
	"github.com/caixw/blogit/internal/vars"
)

func TestBuild(t *testing.T) {
	a := assert.New(t)

	// Dir
	destDir, err := testdata.Temp()
	a.NotError(err)
	dest := builder.DirFS(destDir)
	a.NotError(Build(testdata.Source, dest))
	a.True(filesystem.Exists(dest, "index"+vars.Ext)).
		True(filesystem.Exists(dest, "tags"+vars.Ext)).
		True(filesystem.Exists(dest, "tags/default"+vars.Ext)).
		True(filesystem.Exists(dest, "posts/p1"+vars.Ext))

	// Memory
	dest = builder.MemoryFS()
	a.NotError(Build(testdata.Source, dest))
	a.True(filesystem.Exists(dest, "index"+vars.Ext)).
		True(filesystem.Exists(dest, "tags"+vars.Ext)).
		True(filesystem.Exists(dest, "tags/default"+vars.Ext)).
		True(filesystem.Exists(dest, "posts/p1"+vars.Ext))
}
