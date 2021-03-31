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

	"github.com/caixw/blogit/builder"
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

const postUsage = `创建新文章

执行该命令会在 posts 目录下创建一个同名的文件，如果未指定扩展名，
则自动添加 .md 作为扩展名，可以带目录结构，比如：blogit post 2021/03/31/file
会在项目上目录下添加 posts/2021/03/31.file.md 文件，并在文件内添加必要的字段内容。

执行命令时，当前工作目录必须为项目的根目录。
`

// InitPost 注册 post 子命令
func InitPost(opt *cmdopt.CmdOpt, succ, erro *console.Logger) {
	postFS = opt.New("post", postUsage, post(succ, erro))
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

func getWD() (builder.WritableFS, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return builder.DirFS(dir), nil
}
