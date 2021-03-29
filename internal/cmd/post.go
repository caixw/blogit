// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/issue9/cmdopt"

	"github.com/caixw/blogit/internal/vars"
)

const content = `---
title: title
created: %s
tags:
  - default
state: draft
---

此处书写文章的具体内容
`

var postFS *flag.FlagSet

// initPost 注册 post 子命令
func initPost(opt *cmdopt.CmdOpt) {
	postFS = opt.New("post", "创建新文章\n", post)
}

func post(w io.Writer) error {
	if postFS.NArg() != 1 {
		erro.println("必须指定路径")
		return nil
	}

	wfs, err := getWD()
	if err != nil {
		erro.printf(err.Error())
		return nil
	}

	p := postFS.Arg(0)
	if strings.ToLower(path.Ext(p)) != vars.MarkdownExt {
		p += vars.MarkdownExt
	}
	p = path.Clean(path.Join(vars.PostsDir, p))

	c := fmt.Sprintf(content, time.Now().Format(time.RFC3339))
	if err := wfs.WriteFile(p, []byte(c), os.ModePerm); err != nil {
		erro.println(err)
		return nil
	}
	succ.printf("创建文件: %s\n", p)

	return nil
}
