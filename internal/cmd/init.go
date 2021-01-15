// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
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

// initInit 注册 init 子命令
func initInit(opt *cmdopt.CmdOpt) {
	initFS = opt.New("init", "初始化新的博客内容", initF)
}

func initF(w io.Writer) error {
	if initFS.NArg() != 1 {
		erro.println("必须指定目录")
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		erro.println(err.Error())
		return nil
	}

	// conf.yaml
	conf := &loader.Config{
		Title:  "example",
		URL:    "https://example.com",
		Uptime: time.Now(),
		Theme:  "default",
	}
	path := filepath.Join(dir, vars.ConfYAML)
	if err := writeYAML(path, conf); err != nil {
		erro.println(err.Error())
		return nil
	}
	succ.printf("创建了文件: %s", path)

	// tags.yaml
	tags := []*loader.Tag{
		{
			Slug:    "default",
			Title:   "默认",
			Content: "这是默认的标签",
		},
	}
	path = filepath.Join(dir, vars.TagsYAML)
	if err := writeYAML(path, tags); err != nil {
		erro.println(err.Error())
		return nil
	}
	succ.printf("创建了文件: %s", path)

	// themes
	if err := os.MkdirAll(filepath.Join(dir, vars.ThemesDir), os.ModePerm); err != nil {
		erro.println(err.Error())
		return nil
	}

	if err := os.MkdirAll(filepath.Join(dir, vars.PostsDir), os.ModePerm); err != nil {
		erro.println(err.Error())
	}

	return nil
}

func writeYAML(path string, v interface{}) error {
	bs, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, os.ModePerm)
}
