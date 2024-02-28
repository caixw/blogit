// SPDX-FileCopyrightText: 2020-2024 caixw
//
// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2/styles"
	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"golang.org/x/text/message"
)

const stylesUsage = localeutil.StringPhrase("styles usage")

func initStyles(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("styles", stylesUsage.LocaleString(p), stylesUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			_, err := fmt.Fprintln(w, styles.Names())
			return err
		}
	})
}
