// SPDX-License-Identifier: MIT

package create

import (
	"os"
	"path"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit/filesystem"
	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/vars"
)

func TestCmd_Post(t *testing.T) {
	a := assert.New(t)
	opt := &cmdopt.CmdOpt{}
	succ := &console.Logger{Out: os.Stdout}
	erro := &console.Logger{Out: os.Stderr}
	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)
	a.NotError(os.Chdir(dir))

	p := "2010/01/p1.md"
	InitPost(opt, succ, erro)
	a.NotError(opt.Exec([]string{"post", p}))

	fs := os.DirFS(dir)
	a.True(filesystem.Exists(fs, path.Join(vars.PostsDir, p)))
}
