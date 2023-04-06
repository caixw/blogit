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

	"github.com/caixw/blogit/v2"
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

// InitPost 注册 post 子命令
func InitPost(opt *cmdopt.CmdOpt, succ, erro *console.Logger, lp *message.Printer) {
	opt.New("post", lp.Sprintf("init post usage"), lp.Sprintf("init post usage"), func(fs *flag.FlagSet) cmdopt.DoFunc {
		return func(w io.Writer) error {
			if fs.NArg() != 1 {
				erro.Println(lp.Printf("miss argument"))
				return nil
			}

			wfs, err := getWD()
			if err != nil {
				erro.Println(err)
				return nil
			}

			p := fs.Arg(0)
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
			succ.Println(lp.Sprintf("create file", p))

			return nil
		}
	})
}

func getWD() (blogit.WritableFS, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return blogit.DirFS(dir), nil
}
