// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/text/message"

	"github.com/issue9/cmdopt"
)

func initStyles(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("styles", p.Sprintf("styles usage"), p.Sprintf("styles usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			_, err := fmt.Fprintln(w, styles.Names())
			return err
		}
	})
}
