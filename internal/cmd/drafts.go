// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/internal/data"
	"github.com/caixw/blogit/v2/internal/vars"
)

const (
	draftsTitle    = localeutil.StringPhrase("drafts title")
	draftsUsage    = localeutil.StringPhrase("drafts usage")
	draftsSrcUsage = localeutil.StringPhrase("drafts src usage")
)

func initDrafts(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("drafts", draftsTitle.LocaleString(p), draftsUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		var draftsSrc string
		fs.StringVar(&draftsSrc, "src", "./", draftsSrcUsage.LocaleString(p))

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
