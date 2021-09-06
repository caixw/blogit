// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"io"
	"io/fs"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/testdata"
)

var initFS *flag.FlagSet

// InitInit 注册 init 子命令
func InitInit(opt *cmdopt.CmdOpt, erro *console.Logger, p *message.Printer) {
	initFS = opt.New("init", p.Sprintf("init usage"), initF(erro, p))
}

func initF(erro *console.Logger, p *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if initFS.NArg() != 1 {
			erro.Println(p.Sprintf("miss argument"))
			return nil
		}

		wfs := blogit.DirFS(initFS.Arg(0))
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
