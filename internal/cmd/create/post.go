// SPDX-License-Identifier: MIT

package create

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/issue9/cmdopt"
	"golang.org/x/text/message"

	"github.com/caixw/blogit/v2/builder"
	"github.com/caixw/blogit/v2/internal/cmd/console"
	"github.com/caixw/blogit/v2/internal/vars"
)

const postContent = `---
title: title
created: %s
modified: %s
tags:
  - default
state: draft
---

此处书写文章的具体内容
`

var postFS *flag.FlagSet

// InitPost 注册 post 子命令
func InitPost(opt *cmdopt.CmdOpt, succ, erro *console.Logger, p *message.Printer) {
	postFS = opt.New("post", p.Sprintf("init post usage"), post(succ, erro, p))
}

func post(succ, erro *console.Logger, localePrinter *message.Printer) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if postFS.NArg() != 1 {
			erro.Println(localePrinter.Printf("miss argument"))
			return nil
		}

		wfs, err := getWD()
		if err != nil {
			erro.Println(err)
			return nil
		}

		p := postFS.Arg(0)
		if strings.ToLower(path.Ext(p)) != vars.MarkdownExt {
			p += vars.MarkdownExt
		}
		p = path.Clean(path.Join(vars.PostsDir, p))

		now := time.Now().Format(time.RFC3339)
		c := fmt.Sprintf(postContent, now, now)
		if err := wfs.WriteFile(p, []byte(c), os.ModePerm); err != nil {
			erro.Println(err)
			return nil
		}
		succ.Println(localePrinter.Sprintf("create file", p))

		return nil
	}
}

func getWD() (builder.WritableFS, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return builder.DirFS(dir), nil
}
