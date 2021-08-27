// SPDX-License-Identifier: MIT

package blogit

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"

	"github.com/caixw/blogit/v2/internal/filesystem"
	"github.com/caixw/blogit/v2/internal/testdata"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(version.SemVerValid(Version(true)))
	a.True(version.SemVerValid(Version(false)))
}

func TestBuild(t *testing.T) {
	a := assert.New(t)

	// Dir
	destDir, err := testdata.Temp()
	a.NotError(err)
	dest := DirFS(destDir)
	a.NotError(Build(testdata.Source, dest, nil))
	a.True(filesystem.Exists(dest, "index"+vars.Ext)).
		True(filesystem.Exists(dest, "tags"+vars.Ext)).
		True(filesystem.Exists(dest, "tags/default"+vars.Ext)).
		True(filesystem.Exists(dest, "posts/p1"+vars.Ext))

	// Memory
	dest = MemoryFS()
	a.NotError(Build(testdata.Source, dest, nil))
	a.True(filesystem.Exists(dest, "index"+vars.Ext)).
		True(filesystem.Exists(dest, "tags"+vars.Ext)).
		True(filesystem.Exists(dest, "tags/default"+vars.Ext)).
		True(filesystem.Exists(dest, "posts/p1"+vars.Ext))
}
