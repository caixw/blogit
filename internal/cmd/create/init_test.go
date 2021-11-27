// SPDX-License-Identifier: MIT

package create

import (
	"os"
	"path"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/cmdopt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/filesystem"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestCmd_Init(t *testing.T) {
	a := assert.New(t, false)
	opt := &cmdopt.CmdOpt{}
	erro := &console.Logger{Out: os.Stderr}
	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)

	InitInit(opt, erro, message.NewPrinter(language.Chinese))
	a.NotError(opt.Exec([]string{"init", dir}))

	fs := os.DirFS(dir)
	a.True(filesystem.Exists(fs, vars.ConfYAML)).
		True(filesystem.Exists(fs, vars.TagsYAML)).
		True(filesystem.Exists(fs, path.Join(vars.ThemesDir, "default", vars.ThemeYAML))).
		True(filesystem.Exists(fs, path.Join(vars.PostsDir, "2020/p2.md")))
}
