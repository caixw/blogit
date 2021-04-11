// SPDX-License-Identifier: MIT

package create

import (
	"os"
	"path"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/filesystem"
	"github.com/caixw/blogit/internal/vars"
)

func TestCmd_Init(t *testing.T) {
	a := assert.New(t)
	opt := &cmdopt.CmdOpt{}
	erro := &console.Logger{Out: os.Stderr}
	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)

	InitInit(opt, erro)
	a.NotError(opt.Exec([]string{"init", dir}))

	fs := os.DirFS(dir)
	a.True(filesystem.Exists(fs, vars.ConfYAML)).
		True(filesystem.Exists(fs, vars.TagsYAML)).
		True(filesystem.Exists(fs, path.Join(vars.ThemesDir, "default", vars.ThemeYAML))).
		True(filesystem.Exists(fs, path.Join(vars.PostsDir, "2020/p2.md")))
}
