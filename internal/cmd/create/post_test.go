// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"os"
	"path"
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v3/colors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/filesystem"
	"github.com/caixw/blogit/v2/internal/vars"
)

func TestCmd_Post(t *testing.T) {
	a := assert.New(t, false)
	opt := cmdopt.New(os.Stdout, flag.ContinueOnError, "", nil, nil)
	succ := &console.Logger{Colorize: colors.New(os.Stdout)}
	erro := &console.Logger{Colorize: colors.New(os.Stderr)}
	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)
	a.NotError(os.Chdir(dir))

	p := "2010/01/p1.md"
	InitPost(opt, succ, erro, message.NewPrinter(language.Chinese))
	a.NotError(opt.Exec([]string{"post", p}))

	fs := os.DirFS(dir)
	a.True(filesystem.Exists(fs, path.Join(vars.PostsDir, p)))
}
