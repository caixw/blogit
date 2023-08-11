// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"io"
	"io/fs"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2"
	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/testdata"
)

const initUsage = localeutil.StringPhrase("init usage")

// InitInit 注册 init 子命令
func InitInit(opt *cmdopt.CmdOpt, erro *console.Logger, p *message.Printer) {
	opt.New("init", initUsage.LocaleString(p), initUsage.LocaleString(p), func(initFS *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			if initFS.NArg() != 1 {
				erro.Println(localeutil.StringPhrase("miss argument").LocaleString(p))
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
	})
}
