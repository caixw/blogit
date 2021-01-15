// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	postFS = opt.New("post", "创建新文章", post)
}

func post(w io.Writer) error {
	if postFS.NArg() != 1 {
		erro.println("必须指定路径")
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		erro.printf(err.Error())
		return nil
	}
	dir = filepath.Join(dir, vars.PostsDir)

	path := postFS.Arg(0)
	if strings.ToLower(filepath.Ext(path)) != ".md" {
		path += ".md"
	}
	path = filepath.Clean(filepath.Join(dir, path))

	if !strings.HasPrefix(path, dir) {
		erro.printf("必须位于 %s 目录下\n", dir)
		return nil
	}

	c := fmt.Sprintf(content, time.Now().Format(time.RFC3339))
	if err := ioutil.WriteFile(path, []byte(c), os.ModePerm); err != nil {
		erro.println(err.Error())
		return nil
	}
	succ.printf("创建文件: %s\n", path)

	return nil
}
