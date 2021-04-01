// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"io"
	"io/fs"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit/builder"
	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/testdata"
)

const initUsage = `初始化博客内容

在指定目录下初始化项目的必须文件，比如 conf.yaml、tags.yaml 等文件。
`

var initFS *flag.FlagSet

// InitInit 注册 init 子命令
func InitInit(opt *cmdopt.CmdOpt, succ, erro *console.Logger) {
	initFS = opt.New("init", initUsage, initF(succ, erro))
}

func initF(succ, erro *console.Logger) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if initFS.NArg() != 1 {
			erro.Println("必须指定目录")
			return nil
		}

		wfs := builder.DirFS(initFS.Arg(0))
		return fs.WalkDir(testdata.Source, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			data, err := fs.ReadFile(testdata.Source, path)
			if err != nil {
				return err
			}
			return wfs.WriteFile(path, data, fs.ModePerm)
		})
	}
}
