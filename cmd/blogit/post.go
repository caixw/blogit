// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caixw/blogit/internal/vars"
	"github.com/issue9/cmdopt"
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

func initPost(opt *cmdopt.CmdOpt) {
	postFS = opt.New("post", "创建新文章", post)
}

func post(w io.Writer) error {
	if initFS.NArg() != 1 {
		return errors.New("必须指定路径")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	dir = filepath.Join(dir, vars.PostsDir)

	path := initFS.Arg(0)
	if strings.ToLower(filepath.Ext(path)) != ".md" {
		path += ".md"
	}
	path = filepath.Clean(filepath.Join(dir, path))

	if !strings.HasPrefix(path, dir) {
		return fmt.Errorf("必须位于 %s 目录下", dir)
	}

	c := fmt.Sprintf(content, time.Now().Format(time.RFC3339))
	return ioutil.WriteFile(path, []byte(c), os.ModePerm)
}
