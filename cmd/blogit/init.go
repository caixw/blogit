// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/issue9/cmdopt"
	"gopkg.in/yaml.v2"

	"github.com/caixw/blogit/internal/loader"
	"github.com/caixw/blogit/internal/vars"
)

var initFS *flag.FlagSet

func initInit(opt *cmdopt.CmdOpt) {
	initFS = opt.New("init", "初始化新的博客内容", initF)
}

func initF(w io.Writer) error {
	if initFS.NArg() != 1 {
		return fmt.Errorf("必须指定目录")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	conf := &loader.Config{
		Title:  "example",
		URL:    "https://example.com",
		Uptime: time.Now(),
		Theme:  "default",
	}
	if err := writeYAML(filepath.Join(dir, vars.ConfYAML), conf); err != nil {
		return err
	}

	tags := []*loader.Tag{
		{
			Slug:    "default",
			Title:   "默认",
			Content: "这是默认的标签",
		},
	}
	if err := writeYAML(filepath.Join(dir, vars.TagsYAML), tags); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(dir, vars.ThemesDir), os.ModePerm); err != nil {
		return err
	}

	return os.MkdirAll(filepath.Join(dir, vars.PostsDir), os.ModePerm)
}

func writeYAML(path string, v interface{}) error {
	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, os.ModePerm)
}
