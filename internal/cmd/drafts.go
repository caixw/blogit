// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/vars"
)

func initDrafts(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("drafts", p.Sprintf("drafts title"), p.Sprintf("drafts usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var draftsSrc string
		fs.StringVar(&draftsSrc, "src", "./", p.Sprintf("drafts src"))

		return func(w io.Writer) error {
			d, err := data.Load(os.DirFS(draftsSrc), true, "")
			if err != nil {
				return err
			}

			for _, p := range d.Posts {
				draft := strings.HasPrefix(p.Title, vars.DraftTitleAround) && strings.HasSuffix(p.Title, vars.DraftTitleAround)
				if draft {
					fmt.Fprintln(w, p.Title, "\t", p.Slug)
				}
			}

			return nil
		}
	})
}
