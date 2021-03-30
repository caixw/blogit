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

	"github.com/caixw/blogit/internal/cmd/console"
	"github.com/caixw/blogit/internal/vars"
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
func InitPost(opt *cmdopt.CmdOpt, succ, erro *console.Logger) {
	postFS = opt.New("post", "创建新文章\n", post(succ, erro))
}

func post(succ, erro *console.Logger) cmdopt.DoFunc {
	return func(w io.Writer) error {
		if postFS.NArg() != 1 {
			erro.Println("必须指定路径")
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
		succ.Println("创建文件:", p)

		return nil
	}
}
