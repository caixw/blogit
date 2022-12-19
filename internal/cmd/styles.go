// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2/styles"
	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"
)

func initStyles(opt *cmdopt.CmdOpt, p *message.Printer) {
	opt.New("styles", p.Sprintf("styles usage"), printStyles)
}

func printStyles(w io.Writer) error {
	_, err := fmt.Fprintln(w, styles.Names())
	return err
}
