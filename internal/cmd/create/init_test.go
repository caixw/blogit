// SPDX-License-Identifier: MIT

package create

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/issue9/assert"
	"github.com/issue9/cmdopt"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/builder"
	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/filesystem"
	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

func TestWriteYAML(t *testing.T) {
	a := assert.New(t)

	obj := &loader.Theme{Description: "desc", URL: "https://example.com"}

	wfs := builder.MemoryFS()
	a.NotError(writeYAML(wfs, "conf.yaml", obj))

	f, err := wfs.Open("conf.yaml")
	a.NotError(err).NotNil(f)
	data, err := io.ReadAll(f)
	a.NotError(err).NotNil(data)

	inst := &loader.Theme{}
	a.NotError(yaml.Unmarshal(data, inst))
	a.Equal(inst, obj)
}

func TestCmd_Init(t *testing.T) {
	a := assert.New(t)
	opt := &cmdopt.CmdOpt{}
	succ := &console.Logger{Out: os.Stdout}
	erro := &console.Logger{Out: os.Stderr}
	dir, err := os.MkdirTemp(os.TempDir(), "blogit")
	a.NotError(err)

	InitInit(opt, succ, erro)
	a.NotError(opt.Exec([]string{"init", dir}))

	fs := os.DirFS(dir)
	a.True(filesystem.Exists(fs, vars.ConfYAML)).
		True(filesystem.Exists(fs, vars.TagsYAML)).
		True(filesystem.Exists(fs, path.Join(vars.ThemesDir, "default", vars.ThemeYAML))).
		True(filesystem.Exists(fs, path.Join(vars.PostsDir, time.Now().Format("2006"), "post1.md")))
}
